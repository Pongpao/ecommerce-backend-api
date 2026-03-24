package services

import (
	"errors"
	"project-e-commerce/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CanTransition(current, next string) bool {
	allowed := map[string][]string{
		models.OrderPending: {
			models.OrderPaid,
			models.OrderCancelled,
		},
		models.OrderPaid: {
			models.OrderProcessing,
			models.OrderCancelled,
			models.OrderRefunded,
		},
		models.OrderProcessing: {
			models.OrderShipped,
			models.OrderRefunded,
		},
		models.OrderShipped: {
			models.OrderCompleted,
		},
	}

	for _, status := range allowed[current] {
		if status == next {
			return true
		}
	}

	return false
}
func ChangeOrderStatus(tx *gorm.DB, order *models.Order, newStatus string, changedBy uuid.UUID) error {

	if !CanTransition(order.Status, newStatus) {
		return errors.New("invalid status transition")
	}

	history := models.OrderStatusHistory{
		OrderID:    order.ID,
		FromStatus: order.Status,
		ToStatus:   newStatus,
		ChangedBy:  changedBy,
	}

	order.Status = newStatus

	if err := tx.Save(order).Error; err != nil {
		return err
	}

	if err := tx.Create(&history).Error; err != nil {
		return err
	}

	return nil
}