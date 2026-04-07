package service

import (
	"order-service/internal/model"
	"order-service/internal/repository"
)

type OrderService struct {
	orderRepository *repository.OrderRepository
}

func NewOrderService(orderRepository *repository.OrderRepository) *OrderService {
	return &OrderService{orderRepository: orderRepository}
}

func (s *OrderService) Create(userID uint, productName string) (*model.Order, error) {
	order := &model.Order{
		UserID:      userID,
		ProductName: productName,
	}

	if err := s.orderRepository.Create(order); err != nil {
		return nil, err
	}

	return order, nil
}

func (s *OrderService) List(userID uint) ([]model.Order, error) {
	return s.orderRepository.ListByUserID(userID)
}

func (s *OrderService) Delete(orderID, userID uint) error {
	return s.orderRepository.DeleteByIDAndUserID(orderID, userID)
}
