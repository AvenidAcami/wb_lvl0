package repository

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gomodule/redigo/redis"
	"gorm.io/gorm"
	"log"
	"strings"
	"time"
	"wb_lvl0/internal/model"
)

type OrderRepository struct {
	db  *gorm.DB
	rdb *redis.Pool
}

func NewOrderRepository(db *gorm.DB, rdb *redis.Pool) *OrderRepository {
	return &OrderRepository{db: db,
		rdb: rdb}
}

type IOrderRepository interface {
	InsertOrder(order model.Order) error
	GetOrder(orderUid string, ctx context.Context) (model.Order, error)
}

func (repo *OrderRepository) insertOrderTransaction(ctx context.Context, order model.Order) error {
	var itemsToInsert []model.ItemDB
	tx := repo.db.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := tx.Error; err != nil {
		return err
	}

	// TODO: Вынести перенос данных структуры для вставки в отдельные функции

	// Вставка данных о payment
	paymentToInsert := model.PaymentDB{
		Payment: order.Payment,
	}
	if err := tx.Table("payments").Create(&paymentToInsert).Error; err != nil {
		tx.Rollback()
		return errors.New("something wrong with data in payment section")
	}

	// Вставка данных о delivery
	deliveryToInsert := model.DeliveryDB{
		Delivery: order.Delivery,
	}
	if err := tx.Table("delivery_params").Create(&deliveryToInsert).Error; err != nil {
		tx.Rollback()
		return errors.New("something wrong with data in delivery section")
	}

	// Вставка данных о order
	orderToInsert := model.OrderDB{
		OrderUid:          order.OrderUid,
		TrackNumber:       order.TrackNumber,
		Entry:             order.Entry,
		DeliveryParamsId:  deliveryToInsert.Id,
		PaymentId:         paymentToInsert.Id,
		Locale:            order.Locale,
		InternalSignature: order.InternalSignature,
		CustomerId:        order.CustomerId,
		DeliveryService:   order.DeliveryService,
		Shardkey:          order.Shardkey,
		SmId:              order.SmId,
		DateCreated:       order.DateCreated,
		OofShard:          order.OofShard,
	}
	if err := tx.Table("orders").Create(&orderToInsert).Error; err != nil {
		tx.Rollback()
		return errors.New("something wrong with data in order section")
	}

	// Вставка данных о items
	for _, it := range order.Items {
		itemsToInsert = append(itemsToInsert, model.ItemDB{
			Item:     it,
			OrderUid: order.OrderUid,
		})
	}
	if err := tx.Table("ordered_items").Create(&itemsToInsert).Error; err != nil {
		tx.Rollback()
		return errors.New("something wrong with data in item section")
	}

	return tx.Commit().Error
}

func (repo *OrderRepository) InsertOrder(order model.Order) error {
	var err error
	retries := 5

	conn := repo.rdb.Get()
	defer conn.Close()

	for i := 1; i <= retries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		err = repo.insertOrderTransaction(ctx, order)
		cancel()
		if err == nil {
			redisOrder, err := json.Marshal(&order)
			if err == nil {

				_, err = conn.Do("SETEX", "order:"+order.OrderUid, 300, redisOrder)
				if err != nil {
					log.Println(err)
				}
			} else {
				log.Println(err)
			}
			return nil
		}

		if strings.Contains(err.Error(), "something wrong with data") {
			return err
		}

		if strings.Contains(err.Error(), "duplicate key value violates unique constraint \"orders_pkey\"") {
			return errors.New("order already exists")
		}

		if i < retries {
			time.Sleep(1 * time.Second)
		}

	}

	return errors.New("something went wrong")
}

func (repo *OrderRepository) GetOrder(orderUid string, ctx context.Context) (model.Order, error) {
	var order model.Order
	var orderDB model.OrderDB
	var paymentDB model.PaymentDB
	var deliveryDB model.DeliveryDB
	var itemsDB []model.ItemDB

	conn := repo.rdb.Get()
	defer conn.Close()

	data, err := redis.Bytes(conn.Do("GET", "order:"+orderUid))
	if err == nil {
		err = json.Unmarshal(data, &order)
		if err == nil {
			return order, nil
		} else {
			log.Println(err)
		}
	} else {
		log.Println(err)
	}

	err = repo.db.WithContext(ctx).Table("orders").Where("order_uid = ?", orderUid).First(&orderDB).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return order, errors.New("order not found")
		}
		if strings.Contains(err.Error(), "invalid input syntax for type uuid") {
			return order, errors.New("invalid UUID format")
		}
		return order, err
	}

	err = repo.db.WithContext(ctx).Table("payments").Where("id = ?", orderDB.PaymentId).First(&paymentDB).Error
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return order, errors.New("order not found")
		}
		return order, err
	}

	err = repo.db.WithContext(ctx).Table("delivery_params").Where("id = ?", orderDB.DeliveryParamsId).First(&deliveryDB).Error
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return order, errors.New("order not found")
		}
		return order, err
	}

	err = repo.db.WithContext(ctx).Table("ordered_items").Where("order_uid = ?", orderDB.OrderUid).Find(&itemsDB).Error
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return order, errors.New("order not found")
		}
		return order, err
	}

	order = repo.convertOrdersDbToOrder(orderDB, itemsDB, deliveryDB, paymentDB)

	redisOrder, err := json.Marshal(&order)
	if err == nil {
		_, err = conn.Do("SETEX", "order:"+orderUid, 300, redisOrder)
		if err != nil {
			log.Println(err)
		}
	} else {
		log.Println(err)
	}

	return order, nil
}

func (repo *OrderRepository) convertOrdersDbToOrder(orderDB model.OrderDB, itemsDB []model.ItemDB, deliveryDB model.DeliveryDB, paymentDB model.PaymentDB) model.Order {
	var order model.Order
	// Заполнение основных полей
	order.OrderUid = orderDB.OrderUid
	order.TrackNumber = orderDB.TrackNumber
	order.Entry = orderDB.Entry
	order.Locale = orderDB.Locale
	order.InternalSignature = orderDB.InternalSignature
	order.CustomerId = orderDB.CustomerId
	order.DeliveryService = orderDB.DeliveryService
	order.Shardkey = orderDB.Shardkey
	order.SmId = orderDB.SmId
	order.DateCreated = orderDB.DateCreated
	order.OofShard = orderDB.OofShard

	// Заполнение вложенных полей
	order.Items = repo.convertItemsDbToItems(itemsDB)
	order.Delivery = deliveryDB.Delivery
	order.Payment = paymentDB.Payment

	return order
}

func (repo *OrderRepository) convertItemsDbToItems(itemsDb []model.ItemDB) []model.Item {
	var items []model.Item
	for _, item := range itemsDb {
		items = append(items, item.Item)
	}
	return items
}
