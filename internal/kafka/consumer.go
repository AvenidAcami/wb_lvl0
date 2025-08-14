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
	// Бесконечный цикл получения заказов от кафки
	for {
		ctx := context.Background()
		m, err := reader.ReadMessage(ctx)
		if err != nil {
			log.Println("Error reading message:", err)
			continue
		}
		log.Println("Message received:", string(m.Value))
		var order model.Order
		err = json.Unmarshal(m.Value, &order)
		if err != nil {
			log.Println("Error unmarshalling message:", err)
			continue
		}

		err = orderService.CreateOrder(order)
		if err != nil {
			log.Println("Error creating order:", err)
			continue
		}
		err = reader.CommitMessages(ctx, m)
		if err != nil {
			log.Println("Error committing message:", err)
		}
	}
}
