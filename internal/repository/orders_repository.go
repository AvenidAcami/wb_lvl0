package repository

import (
	"gorm.io/gorm"
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
	return nil
}
