package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
	"time"
	"wb_lvl0/internal/model"
	"wb_lvl0/internal/repository"
)

type OrderService struct {
	repo repository.IOrderRepository
	rdb  *redis.Pool
}

func NewOrderService(repo repository.IOrderRepository, rdb *redis.Pool) *OrderService {
	return &OrderService{repo: repo, rdb: rdb}
}

type IOrderService interface {
	CreateOrder(order model.Order) error
	GetOrder(ctx context.Context, orderUid string) (model.Order, error)
	RestoreCache() error
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
			return fmt.Errorf("%s is required", name)

		}
	}
	return nil
}

func (serv *OrderService) CreateOrder(order model.Order) error {
	err := serv.validateOrderInfo(order)
	if err != nil {
		return err
	}

	conn := serv.rdb.Get()
	defer conn.Close()
	serv.setOrderInRedis(conn, order)

	return serv.repo.InsertOrder(order)
}

func (serv *OrderService) GetOrder(ctx context.Context, orderUid string) (model.Order, error) {
	timeoutCtx, cancel := serv.createContext(ctx, 3)
	defer cancel()

	conn := serv.rdb.Get()
	defer conn.Close()

	order, err := serv.getOrderFromRedis(conn, orderUid)
	if err == nil {
		return order, nil
	} else {
		log.Println("getOrderFromRedis error:", err)
	}

	order, err = serv.repo.GetOrder(orderUid, timeoutCtx)
	if err != nil {
		log.Println(err)
		return model.Order{}, err
	}

	serv.setOrderInRedis(conn, order)
	return order, nil
}

func (serv *OrderService) createContext(ctx context.Context, seconds time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, seconds*time.Second)
}

func (serv *OrderService) RestoreCache() error {
	ctx, cancel := serv.createContext(context.Background(), 5)
	defer cancel()
	orders, err := serv.repo.GetLastOrders(ctx)
	if err != nil {
		return err
	}

	conn := serv.rdb.Get()
	defer conn.Close()

	for _, order := range orders {
		orderstr, err := json.Marshal(order)
		if err != nil {
			continue
		}
		_, _ = conn.Do("SETEX", "order:"+order.OrderUid, 300, orderstr)
	}
	return nil
}

func (serv *OrderService) setOrderInRedis(conn redis.Conn, order model.Order) {

	redisOrder, err := json.Marshal(&order)
	if err == nil {
		_, err = conn.Do("SETEX", "order:"+order.OrderUid, 300, redisOrder)
		if err != nil {
			log.Println(err)
		}
	} else {
		log.Println(err)
	}
}

func (serv *OrderService) getOrderFromRedis(conn redis.Conn, orderUid string) (model.Order, error) {
	var order model.Order

	data, err := redis.Bytes(conn.Do("GET", "order:"+orderUid))
	if err == nil {
		err = json.Unmarshal(data, &order)
		if err == nil {
			return order, nil
		} else {
			return model.Order{}, err
		}
	} else {
		log.Println(err)
	}
	return order, err
}
