package app

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"time"
	"wb_lvl0/config"
	"wb_lvl0/internal/controller"
	"wb_lvl0/internal/kafka"
	"wb_lvl0/internal/repository"
	"wb_lvl0/internal/router"
	"wb_lvl0/internal/service"
)

func Run() {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Инициализация компонентов
	config.InitENV()
	db := config.InitPostgres()
	redisPool := config.InitRedis()
	reader := config.InitKafka()

	orderRepository := repository.NewOrderRepository(db)
	orderService := service.NewOrderService(orderRepository, redisPool)
	orderController := controller.NewOrdersController(orderService)
	router.InitOrderRoutes(r, orderController)

	// Восстановление кэша (1000 последних записей)
	err := orderService.RestoreCache()
	if err != nil {
		log.Println("Error restoring cache")
	}

	// Запуск консьюмера
	go kafka.ParseOrders(reader, orderService)

	r.Run(":8081")
}
