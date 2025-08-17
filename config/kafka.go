package config

import (
	"github.com/segmentio/kafka-go"
	"log"
	"time"
)

func InitKafka() *kafka.Reader {
	for {
		conn, err := kafka.Dial("tcp", "kafka:9092")
		if err == nil {
			conn.Close()
			log.Println(err)
			break
		}
		time.Sleep(5 * time.Second)
	}

	for {
		conn, err := kafka.Dial("tcp", "kafka:9092")
		if err != nil {
			log.Println(err)
			time.Sleep(5 * time.Second)
			continue
		}

		partitions, err := conn.ReadPartitions("orders")
		conn.Close()
		if err == nil && len(partitions) > 0 {
			break
		}

		time.Sleep(5 * time.Second)
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"kafka:9092"},
		Topic:    "orders",
		GroupID:  "order-consumers",
		MaxBytes: 10e6,
	})

	log.Println("Kafka reader initialized")
	return reader
}
