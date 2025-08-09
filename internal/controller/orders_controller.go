package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wb_lvl0/internal/service"
)

type OrdersController struct {
	serv service.IOrderService
}

func NewOrdersController(serv service.IOrderService) *OrdersController {
	return &OrdersController{serv: serv}
}

func (contr *OrdersController) GetOrder(c *gin.Context) {
	orderUid := c.Param("order_uid")
	order, err := contr.serv.GetOrder(orderUid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"order": order,
	})
}
