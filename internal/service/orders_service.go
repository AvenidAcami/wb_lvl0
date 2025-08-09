package service

import (
	"errors"
	"fmt"
	"wb_lvl0/internal/model"
	"wb_lvl0/internal/repository"
)

type OrderService struct {
	repo repository.IOrderRepository
}

func NewOrderService(repo repository.IOrderRepository) *OrderService {
	return &OrderService{repo: repo}
}

type IOrderService interface {
	CreateOrder(order model.Order) error
	GetOrder(order_uid string) (model.Order, error)
}

func (serv *OrderService) validateOrderInfo(order model.Order) error {
	// Проверка полей order
	err := serv.validateRequiredStrings(map[string]string{
		"order_uid":        order.OrderUid,
		"track_number":     order.TrackNumber,
		"entry":            order.Entry,
		"locale":           order.Locale,
		"customer_id":      order.CustomerId,
		"delivery_service": order.DeliveryService,
		"shardkey":         order.Shardkey,
		"date_created":     order.DateCreated,
		"oof_shard":        order.OofShard,
	})
	if err != nil {
		return err
	}
	err = serv.validateRequiredInts(map[string]int{
		"sm_id": order.SmId,
	})
	if err != nil {
		return err
	}

	// Проверка полей delivery
	err = serv.validateRequiredStrings(map[string]string{
		"delivery.name":    order.Delivery.Name,
		"delivery.phone":   order.Delivery.Phone,
		"delivery.zip":     order.Delivery.Zip,
		"delivery.city":    order.Delivery.City,
		"delivery.address": order.Delivery.Address,
		"delivery.region":  order.Delivery.Region,
		"delivery.email":   order.Delivery.Email,
	})
	if err != nil {
		return err
	}

	// Проверка полей payment
	err = serv.validateRequiredStrings(map[string]string{
		"payment.transaction": order.Payment.Transaction,
		"payment.currency":    order.Payment.Currency,
		"payment.provider":    order.Payment.Provider,
		"payment.bank":        order.Payment.Bank,
	})
	if err != nil {
		return err
	}
	err = serv.validateRequiredInts(map[string]int{
		"payment.amount":      order.Payment.Amount,
		"payment.payment_dt":  order.Payment.PaymentDt,
		"payment.goods_total": order.Payment.GoodsTotal,
	})
	if err != nil {
		return err
	}

	// Проверка полей items
	if len(order.Items) == 0 {
		return errors.New("items must have at least one item")
	}
	for i := 0; i < len(order.Items); i++ {
		err = serv.validateRequiredStrings(map[string]string{
			"item.track_number": order.Items[i].TrackNumber,
			"item.rid":          order.Items[i].Rid,
			"item.name":         order.Items[i].Name,
			"item.size":         order.Items[i].Size,
			"item.brand":        order.Items[i].Brand,
		})
		if err != nil {
			return err
		}

		err = serv.validateRequiredInts(map[string]int{
			"item.chrt_id":     order.Items[i].ChrtId,
			"item.price":       order.Items[i].Price,
			"item.total_price": order.Items[i].TotalPrice,
			"item.nm_id":       order.Items[i].NmId,
			"item.status":      order.Items[i].Status,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (serv *OrderService) validateRequiredStrings(fields map[string]string) error {
	for name, value := range fields {
		if value == "" {
			return fmt.Errorf("%s is required", name)
		}
	}
	return nil
}

func (serv *OrderService) validateRequiredInts(fields map[string]int) error {
	for name, value := range fields {
		if value == 0 {
			return fmt.Errorf("%s is required and must be non-zero", name)
		}
	}
	return nil
}

func (serv *OrderService) CreateOrder(order model.Order) error {
	err := serv.validateOrderInfo(order)
	if err != nil {
		return err
	}
	return serv.repo.InsertOrder(order)
}

func (serv *OrderService) GetOrder(orderUid string) (model.Order, error) {
	return serv.repo.GetOrder(orderUid)
}
