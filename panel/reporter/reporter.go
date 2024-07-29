package reporter

import (
	"YellowBloomKnapsack/mini-yektanet/common/dto"
	"YellowBloomKnapsack/mini-yektanet/common/models"
	"YellowBloomKnapsack/mini-yektanet/panel/database"
	"YellowBloomKnapsack/mini-yektanet/panel/handlers"
	"YellowBloomKnapsack/mini-yektanet/panel/logic"
	"github.com/golang/protobuf/proto"
	"gorm.io/gorm"
	"strconv"
	"time"

	//"context"
	"fmt"
	"log"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type ReporterInterface interface {
	Start()
}

type MessageHandlerFunc func(messages []kafka.Message) error

type ConsumerService struct {
	consumer  *kafka.Consumer
	buffLimit int
	topic     string
	handler   MessageHandlerFunc
}

type ReporterService struct {
	consumers []*ConsumerService
}

func newConsumerService(kafkaBootstrapServers, topic string, buffLimit int, handler MessageHandlerFunc) *ConsumerService {
	kafkaConsumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": kafkaBootstrapServers,
		"auto.offset.reset": "earliest",
		"group.id":          "yektanet-reporter-" + topic,
	})
	if err != nil {
		log.Printf("Failed to create kafka consumer: %v", err)
	}

	err = kafkaConsumer.Subscribe(topic, nil)
	if err != nil {
		log.Printf("Failed to subscribe to topic %s: %v", topic, err)
	}

	return &ConsumerService{
		consumer:  kafkaConsumer,
		buffLimit: buffLimit,
		topic:     topic,
		handler:   handler,
	}
}

func (c *ConsumerService) Start() {
	defer c.consumer.Close()

	buffer := make([]kafka.Message, 0, c.buffLimit)

	for {
		msg, err := c.consumer.ReadMessage(100 * time.Millisecond)
		if err != nil {
			// Handle message reading errors
			//log.Printf("Consumer error: %v\n", err)
			continue
		}

		buffer = append(buffer, *msg)
		if len(buffer) >= c.buffLimit {
			fmt.Println("Buffer full, clearing it")
			if err := c.handler(buffer); err != nil {
				log.Printf("Error handling message: %v", err)
			}
			buffer = make([]kafka.Message, 0) // Reset the buffer
		}
	}

}

func NewReporterService(clickTopic, impressionTopic string, clickBuffLimit, impressionBuffLimit int) ReporterInterface {
	kafkaBootstrapServers := os.Getenv("KAFKA_BOOTSTRAP_SERVERS") + ":9092"
	if kafkaBootstrapServers == "" {
		log.Fatal("KAFKA_BOOTSTRAP_SERVERS environment variable is not set")
	}

	consumers := []*ConsumerService{
		newConsumerService(kafkaBootstrapServers, clickTopic, clickBuffLimit, handleClick),
		newConsumerService(kafkaBootstrapServers, impressionTopic, impressionBuffLimit, handleImpression),
	}

	return &ReporterService{
		consumers: consumers,
	}
}

func (r *ReporterService) Start() {
	for _, consumer := range r.consumers {
		go consumer.Start()
	}
}

func handleClick(messages []kafka.Message) error {
	fmt.Println("Handling click event")

	for _, msg := range messages {
		var event dto.ClickEvent
		err := proto.Unmarshal(msg.Value, &event)
		if err != nil {
			return fmt.Errorf("failed to unmarshal proto message: %w", err)
		}

		// Start a transaction
		tx := database.DB.Begin()

		// Find the publisher
		var publisher models.Publisher
		if err := tx.Where("id = ?", event.PublisherId).First(&publisher).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to get publisher by username: %w", err)
		}

		// Create the interaction
		eventTime, _ := time.Parse(time.RFC3339, event.EventTime)
		interaction := models.AdsInteraction{
			Type:        int(models.Click),
			EventTime:   eventTime,
			AdID:        uint(event.AdId),
			Bid:         event.Bid,
			PublisherID: uint(event.PublisherId),
		}

		if err := tx.Create(&interaction).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to create interaction: %w", err)
		}

		// Find the ad and its associated advertiser
		var ad models.Ad
		if err := tx.Preload("Advertiser").First(&ad, event.AdId).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to get ad by ad id: %w", err)
		}

		// Update ad's total cost
		if err := tx.Model(&ad).Update("total_cost", gorm.Expr("total_cost + ?", ad.Bid)).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to update total_cost: %w", err)
		}

		yektanetPortionString := os.Getenv("YEKTANET_PORTION")

		// Convert the value to an integer
		yektanetPortion, err := strconv.Atoi(yektanetPortionString)
		if err != nil || yektanetPortion < 0 || yektanetPortion > 100 {
			tx.Rollback()
			return fmt.Errorf("failed to convert YEKTANET_PORTION to int: %w", err)
		}

		// Increase publisher's balance
		publisherPortion := ad.Bid * int64(100-yektanetPortion) / 100
		if err := tx.Model(&publisher).Update("balance", gorm.Expr("balance + ?", publisherPortion)).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to update balance: %w", err)
		}

		// Create a new transaction record for advertiser
		transaction_publisher := models.Transaction{
			CustomerID:   uint(event.PublisherId),
			CustomerType: models.Customer_Publisher,
			Amount:       publisherPortion,
			Income:       true,
			Successful:   true,
			Time:         time.Now(),
			Description:  "click on ad",
		}

		if err := tx.Create(&transaction_publisher).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to create transaction publisher: %w", err)
		}

		sum, _ := logic.GetSumOfBids(database.DB, ad.ID)
		if sum > ad.Advertiser.Balance {
			go handlers.NotifyAdsBrake(ad.ID)
		}

		// Decrease advertiser's balance
		if err := tx.Model(&ad.Advertiser).Update("balance", gorm.Expr("balance - ?", ad.Bid)).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("Failed to update advertiser's balance: %w", err)
		}

		// Create a new transaction record for advertiser
		transaction_advertiser := models.Transaction{
			CustomerID:   ad.AdvertiserID,
			CustomerType: models.Customer_Advertiser,
			Amount:       ad.Bid,
			Income:       false,
			Successful:   true,
			Time:         time.Now(),
			Description:  "click on ad",
		}

		if err := tx.Create(&transaction_advertiser).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("Failed to update advertiser's balance: %w", err)
		}

		// Commit the transaction
		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("Failed to commit transaction: %w", err)
		}
	}

	return nil
}

func handleImpression(messages []kafka.Message) error {
	fmt.Println("Handling impression event")
	interactionsToInsert := make([]models.AdsInteraction, 0)
	for _, msg := range messages {
		var event dto.ImpressionEvent
		err := proto.Unmarshal(msg.Value, &event)
		if err != nil {
			return fmt.Errorf("failed to unmarshal proto message: %w", err)
		}

		//// Create the interaction
		eventTime, _ := time.Parse(time.RFC3339, event.EventTime)
		interaction := models.AdsInteraction{
			Type:        int(models.Impression),
			EventTime:   eventTime,
			AdID:        uint(event.AdId),
			PublisherID: uint(event.PublisherId),
			Bid:         event.Bid,
		}

		interactionsToInsert = append(interactionsToInsert, interaction)
	}

	if err := database.DB.Create(&interactionsToInsert).Error; err != nil {
		return err
	}
	fmt.Println("///////////////////??????????????????????????????????5")
	return nil
}
