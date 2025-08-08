package kafka

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"log"
	"wb_lvl0/internal/model"
)

func ParseOrders(reader *kafka.Reader) {
	var order model.Order
	log.Println("Start parse orders")
	for {
		//TODO: Сделать передачу в сервис
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Println("Error reading message:", err)
			continue
		}
		log.Println("Message received:", string(m.Value))
		err = json.Unmarshal(m.Value, &order)
		if err != nil {
			log.Println("Error unmarshalling message:", err)
			continue
		}

	}
}
