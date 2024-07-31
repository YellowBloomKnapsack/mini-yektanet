package reporter

import (
	"YellowBloomKnapsack/mini-yektanet/common/dto"
	"YellowBloomKnapsack/mini-yektanet/common/models"
	"YellowBloomKnapsack/mini-yektanet/panel/database"
	"YellowBloomKnapsack/mini-yektanet/panel/grafana"
	"YellowBloomKnapsack/mini-yektanet/panel/handlers"
	"YellowBloomKnapsack/mini-yektanet/panel/logic"

	// "YellowBloomKnapsack/mini-yektanet/panel/grafana"
	"strconv"
	"time"

	"github.com/golang/protobuf/proto"
	"gorm.io/gorm"

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
	timeout   time.Duration // Add a timeout field for the consumer service
}

type ReporterService struct {
	consumers []*ConsumerService
}

func newConsumerService(kafkaBootstrapServers, topic string, buffLimit int, handler MessageHandlerFunc, timeout time.Duration) *ConsumerService {
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
		timeout:   timeout,
	}
}

func (c *ConsumerService) Start() {
	defer c.consumer.Close()

	buffer := make([]kafka.Message, 0, c.buffLimit)
	timer := time.NewTimer(c.timeout) // Create a new timer with the specified timeout

	for {
		select {
		case <-timer.C:
			// If the timer expires, process the messages in the buffer
			if len(buffer) > 0 {
				log.Println("Timer expired, processing messages")
				if err := c.handler(buffer); err != nil {
					log.Printf("Error handling messages: %v", err)
				}

				buffer = make([]kafka.Message, 0)

			}
			timer.Reset(c.timeout)
		default:
			msg, err := c.consumer.ReadMessage(100 * time.Millisecond)
			if err != nil {
				continue
			}

			buffer = append(buffer, *msg)
			if len(buffer) >= c.buffLimit {
				fmt.Println("Buffer full, processing messages")
				if err := c.handler(buffer); err != nil {
					log.Printf("Error handling messages: %v", err)
				}
				buffer = make([]kafka.Message, 0)
				timer.Reset(c.timeout)
			}
		}
	}
}

func NewReporterService(clickTopic, impressionTopic string, clickBuffLimit, impressionBuffLimit int) ReporterInterface {
	kafkaBootstrapServers := os.Getenv("KAFKA_BOOTSTRAP_SERVERS") + ":9092"
	if kafkaBootstrapServers == "" {
		log.Fatal("KAFKA_BOOTSTRAP_SERVERS environment variable is not set")
	}

	timeoutStr := os.Getenv("TIMEOUT")
	if timeoutStr == "" {
		log.Fatal("TIMEOUT environment variable is not set")
	}

	timeoutInt, err := strconv.Atoi(timeoutStr)
	if err != nil {
		log.Fatalf("Invalid TIMEOUT value: %v", err)
	}

	timeout := time.Duration(timeoutInt) * time.Second

	consumers := []*ConsumerService{
		newConsumerService(kafkaBootstrapServers, clickTopic, clickBuffLimit, handleClick, timeout),
		newConsumerService(kafkaBootstrapServers, impressionTopic, impressionBuffLimit, handleImpression, timeout),
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

		var publisher models.Publisher
		if err := tx.Where("id = ?", event.PublisherId).First(&publisher).Error; err != nil {
			tx.Rollback()
			grafana.TransactionCount.WithLabelValues("click_failure").Inc()
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
			grafana.TransactionCount.WithLabelValues("click_failure").Inc()
			return fmt.Errorf("failed to create interaction: %w", err)
		}

		// Find the ad and its associated advertiser
		var ad models.Ad
		if err := tx.Preload("Advertiser").First(&ad, event.AdId).Error; err != nil {
			tx.Rollback()
			grafana.TransactionCount.WithLabelValues("click_failure").Inc()
			return fmt.Errorf("failed to get ad by ad id: %w", err)
		}

		// Update ad's total cost
		if err := tx.Model(&ad).Update("total_cost", gorm.Expr("total_cost + ?", ad.Bid)).Error; err != nil {
			tx.Rollback()
			grafana.TransactionCount.WithLabelValues("click_failure").Inc()
			return fmt.Errorf("failed to update total_cost: %w", err)
		}

		yektanetPortionString := os.Getenv("YEKTANET_PORTION")

		// Convert the value to an integer
		yektanetPortion, err := strconv.Atoi(yektanetPortionString)
		if err != nil || yektanetPortion < 0 || yektanetPortion > 100 {
			tx.Rollback()
			grafana.TransactionCount.WithLabelValues("click_failure").Inc()
			return fmt.Errorf("failed to convert YEKTANET_PORTION to int: %w", err)
		}

		// Increase publisher's balance
		publisherPortion := ad.Bid * int64(100-yektanetPortion) / 100
		if err := tx.Model(&publisher).Update("balance", gorm.Expr("balance + ?", publisherPortion)).Error; err != nil {
			tx.Rollback()
			grafana.TransactionCount.WithLabelValues("click_failure").Inc()
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
			grafana.TransactionCount.WithLabelValues("click_failure").Inc()
			return fmt.Errorf("failed to create transaction publisher: %w", err)
		}

		sum, _ := logic.GetSumOfBids(database.DB, ad.ID)
		sum += event.Bid
		log.Printf("Ad with id %d has %d cost over the last minute.\n", ad.ID, sum)
		if sum > ad.Advertiser.Balance-event.Bid {
			log.Println("Braking it")
			go handlers.NotifyAdsBrake(ad.ID)
		}

		// Decrease advertiser's balance
		if err := tx.Model(&ad.Advertiser).Update("balance", gorm.Expr("balance - ?", ad.Bid)).Error; err != nil {
			tx.Rollback()
			grafana.TransactionCount.WithLabelValues("click_failure").Inc()
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
			grafana.TransactionCount.WithLabelValues("click_failure").Inc()
			return fmt.Errorf("Failed to update advertiser's balance: %w", err)
		}

		// Commit the transaction
		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			grafana.TransactionCount.WithLabelValues("click_failure").Inc()
			return fmt.Errorf("Failed to commit transaction: %w", err)
		}

		grafana.TotalRevenue.Add(float64(ad.Bid*int64(yektanetPortion)) / 100)
		grafana.TotalPublisherBalance.Add(float64(publisherPortion))
		grafana.TotalAdvertiserBalance.Add(-float64(transaction_advertiser.Amount))
		grafana.ClickCount.Inc()
		grafana.TransactionCount.WithLabelValues("click_success").Inc()
	}

	return nil
}

func handleImpression(messages []kafka.Message) error {
	log.Println("Handling impression event")
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

	grafana.ImpressionCount.Add(float64(len(interactionsToInsert)))
	return nil
}
