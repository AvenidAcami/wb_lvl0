package kafka

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"log"
	"wb_lvl0/internal/model"
	"wb_lvl0/internal/service"
)

func ParseOrders(reader *kafka.Reader, orderService service.IOrderService) {
	var order model.Order
	log.Println("Start parse orders")
	for {
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

		err = orderService.CreateOrder(order)
		if err != nil {
			log.Println("Error creating order:", err)
		}
	}
}
