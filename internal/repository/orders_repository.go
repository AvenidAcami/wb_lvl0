package repository

import (
	"gorm.io/gorm"
	"log"
	"wb_lvl0/internal/model"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

type IOrderRepository interface {
	InsertOrder(order model.Order) error
}

func (repo *OrderRepository) InsertOrder(order model.Order) error {
	var itemsToInsert []model.ItemToInsert
	tx := repo.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	// TODO: Вынести перенос данных структуры для вставки в отдельные функции

	// Вставка данных о payment
	paymentToInsert := model.PaymentToInsert{
		Payment: order.Payment,
	}
	if err := tx.Table("payments").Create(&paymentToInsert).Error; err != nil {
		tx.Rollback()
		log.Println("Ошибка во время вставки payment")
		return err
	}

	// Вставка данных о delivery
	deliveryToInsert := model.DeliveryToInsert{
		Delivery: order.Delivery,
	}
	if err := tx.Table("delivery_params").Create(&deliveryToInsert).Error; err != nil {
		tx.Rollback()
		log.Println("Ошибка во время вставки delivery")
		return err
	}

	// Вставка данных о order
	orderToInsert := model.OrderToInsert{
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
		log.Println("Ошибка во время вставки order")
		return err
	}

	// Вставка данных о items
	for _, it := range order.Items {
		itemsToInsert = append(itemsToInsert, model.ItemToInsert{
			Item:     it,
			OrderUid: order.OrderUid,
		})
	}
	if err := tx.Table("ordered_items").Create(&itemsToInsert).Error; err != nil {
		tx.Rollback()
		log.Println("Ошибка во время вставки items")
		return err
	}

	return tx.Commit().Error
}
