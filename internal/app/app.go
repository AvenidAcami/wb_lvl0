package app

import (
	"github.com/gin-gonic/gin"
	"wb_lvl0/config"
)

func Run() {
	r := gin.Default()
	config.InitENV()
	config.InitPostgres()
	r.Run(":8080")
}
