package repository

import (
	"order-service/internal/model"

	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(order *model.Order) error {
	return r.db.Create(order).Error
}

func (r *OrderRepository) ListByUserID(userID uint) ([]model.Order, error) {
	var orders []model.Order
	if err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&orders).Error; err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *OrderRepository) DeleteByIDAndUserID(orderID, userID uint) error {
	return r.db.Where("id = ? AND user_id = ?", orderID, userID).Delete(&model.Order{}).Error
}
