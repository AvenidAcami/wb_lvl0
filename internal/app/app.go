package app

import (
	"github.com/gin-gonic/gin"
	"log"
	"wb_lvl0/config"
	"wb_lvl0/internal/controller"
	"wb_lvl0/internal/kafka"
	"wb_lvl0/internal/repository"
	"wb_lvl0/internal/router"
	"wb_lvl0/internal/service"
)

func Run() {
	r := gin.Default()
	config.InitENV()
	db := config.InitPostgres()
	redisPool := config.InitRedis()

	orderRepository := repository.NewOrderRepository(db)
	orderService := service.NewOrderService(orderRepository, redisPool)

	err := orderService.RestoreCache()
	if err != nil {
		log.Println("Error restoring cache")
	}
	orderController := controller.NewOrdersController(orderService)

	reader := config.InitKafka()
	go kafka.ParseOrders(reader, orderService)
	router.InitOrderRoutes(r, orderController)
	r.Run(":8080")
}
