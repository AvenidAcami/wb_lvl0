package app

import (
	"github.com/gin-gonic/gin"
	"wb_lvl0/config"
	"wb_lvl0/internal/kafka"
	"wb_lvl0/internal/repository"
	"wb_lvl0/internal/service"
)

func Run() {
	r := gin.Default()
	config.InitENV()
	db := config.InitPostgres()

	orderRepository := repository.NewOrderRepository(db)
	orderService := service.NewOrderService(orderRepository)
	// TODO: Добавить контроллер и инит роутера
	// TODO: Добавить редис и кэширование первого уровня
	// TODO: Добавит тесты

	reader := config.InitKafka()
	go kafka.ParseOrders(reader, orderService)
	r.Run(":8080")
}
