package router

import (
	"github.com/gin-gonic/gin"
	"wb_lvl0/internal/controller"
)

func InitOrderRoutes(r *gin.Engine, contr *controller.OrdersController) {
	orderGroup := r.Group("/order")
	{
		orderGroup.GET(":order_uid", contr.GetOrder)
	}
}
