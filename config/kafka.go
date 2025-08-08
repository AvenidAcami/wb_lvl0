package config

import (
	"github.com/segmentio/kafka-go"
	"log"
)

func InitKafka() *kafka.Reader {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"kafka:9092"},
		Topic:   "orders",
		GroupID: "order-consumers",
	})
	log.Println("Init kafka reader")
	return reader
}
