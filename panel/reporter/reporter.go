package reporter

import (
	//"context"
	"fmt"
	"log"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type ReporterInterface interface {
	Start()
}

type MessageHandlerFunc func(msg kafka.Message) error

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
		"group.id":          1,
	})
	if err != nil {
		log.Fatal(err)
	}

	err = kafkaConsumer.Subscribe(topic, nil)
	if err != nil {
		log.Fatal(err)
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
		msg, err := c.consumer.ReadMessage(-1)
		if err != nil {
			// Handle message reading errors
			log.Printf("Consumer error: %v\n", err)
			continue
		}

		buffer = append(buffer, *msg)
		fmt.Println("HERE")
		if len(buffer) >= c.buffLimit {
			for _, message := range buffer {
				if err := c.handler(message); err != nil {
					log.Printf("Error handling message: %v", err)
				}
			}
			buffer = buffer[:0] // Reset the buffer
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

func handleClick(msg kafka.Message) error {
	fmt.Println("HANDLIGN CLICDKjckdsfjkgfjgklfd")
	return nil
}

func handleImpression(msg kafka.Message) error {
	fmt.Println("handle impression")
	return nil
}

//
//func (kc *ReporterService) HandleMessage(msg kafka.Message) error {
//	topicName := *msg.TopicPartition.Topic
//	if topicName == os.Getenv("KAFKA_TOPIC_CLICK") {
//		var event dto.ClickEvent
//		err := proto.Unmarshal(msg.Value, &event)
//		if err != nil {
//			return fmt.Errorf("failed to unmarshal proto message: %w", err)
//		}
//	} else if topicName == os.Getenv("KAFKA_TOPIC_IMPRESSION") {
//		var event dto.ImpressionEvent
//		err := proto.Unmarshal(msg.Value, &event)
//		if err != nil {
//			return fmt.Errorf("failed to unmarshal proto message: %w", err)
//		}
//	}
//
//	// Log or process the event
//	fmt.Printf("Consumed message from topic %s: %v\n", *msg.TopicPartition.Topic, event)
//	return nil
//}
