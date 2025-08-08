package app

import (
	"github.com/gin-gonic/gin"
	"wb_lvl0/config"
	"wb_lvl0/internal/kafka"
)

func Run() {
	r := gin.Default()
	config.InitENV()
	config.InitPostgres()
	reader := config.InitKafka()
	go kafka.ParseOrders(reader)
	r.Run(":8080")
}
