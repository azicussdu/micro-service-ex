package service

import (
	"log"

	"order-service/internal/events"
	"order-service/internal/model"
	"order-service/internal/repository"
)

type OrderService struct {
	orderRepository *repository.OrderRepository
	publisher       *events.Publisher
}

func NewOrderService(orderRepository *repository.OrderRepository, publisher *events.Publisher) *OrderService {
	return &OrderService{
		orderRepository: orderRepository,
		publisher:       publisher,
	}
}

func (s *OrderService) Create(userID uint, productName string) (*model.Order, error) {
	order := &model.Order{
		UserID:      userID,
		ProductName: productName,
	}

	if err := s.orderRepository.Create(order); err != nil {
		return nil, err
	}

	if s.publisher != nil {
		if err := s.publisher.Publish("order.created", order); err != nil {
			log.Printf("failed to publish order.created event: %v", err)
		}
	}

	return order, nil
}

func (s *OrderService) List(userID uint) ([]model.Order, error) {
	return s.orderRepository.ListByUserID(userID)
}

func (s *OrderService) Delete(orderID, userID uint) error {
	if err := s.orderRepository.DeleteByIDAndUserID(orderID, userID); err != nil {
		return err
	}

	if s.publisher != nil {
		if err := s.publisher.Publish("order.deleted", map[string]uint{
			"order_id": orderID,
			"user_id":  userID,
		}); err != nil {
			log.Printf("failed to publish order.deleted event: %v", err)
		}
	}

	return nil
}
