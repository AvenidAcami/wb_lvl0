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

// GetOrder
// @Summary Get order
// @Tags orders
// @Description Get order by order_uid
// @ID get-order
// @Produce json
// @Param id path string true "Order uid"
// @Success 200 {object} model.Order
// @Failure 400 {object} model.ErrorResponse "Invalid uid format"
// @Failure 404 {object} model.ErrorResponse "Order not found"
// @Failure 504 {object} model.ErrorResponse "Server timeout"
// @Router /order/{id} [get]
func (contr *OrdersController) GetOrder(c *gin.Context) {
	orderUid := c.Param("order_uid")
	order, err := contr.serv.GetOrder(c.Request.Context(), orderUid)
	if err != nil {
		if err.Error() == "context deadline exceeded" {
			c.JSON(http.StatusGatewayTimeout, gin.H{
				"code":  http.StatusGatewayTimeout,
				"error": "server response timed out",
			})
			return
		}
		if err.Error() == "order not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"code":  http.StatusNotFound,
				"error": "order not found",
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  http.StatusBadRequest,
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"order": order,
	})
}
