package controllers

import (
	"fmt"
	"project-e-commerce/config"
	"project-e-commerce/models"
	"project-e-commerce/utils"

	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PayOrder godoc
// @Summary Pay order
// @Tags payment
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} map[string]interface{}
// @Router /orders/{id}/pay [post]
func PayOrder(c *gin.Context) {

	// ✅ ดึง user id
	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.Error(utils.Internal("user id missing in context", nil))
		return
	}

	userIDStr, ok := userIDValue.(string)
	if !ok {
		c.Error(utils.Internal("invalid user id format", nil))
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.Error(utils.Internal("invalid user id", err))
		return
	}

	orderIDParam := c.Param("id")
	orderID, err := uuid.Parse(orderIDParam)
	if err != nil {
		c.Error(utils.BadRequest("invalid order id", err))
		return
	}

	err = config.DB.Transaction(func(tx *gorm.DB) error {

		var order models.Order
		if err := tx.
			Where("id = ? AND user_id = ?", orderID, userID).
			First(&order).Error; err != nil {

			if errors.Is(err, gorm.ErrRecordNotFound) {
				return utils.BadRequest("order not found", nil)
			}

			return utils.Internal("database error", err)
		}

		if order.Status != models.OrderPending {
			return utils.BadRequest("order already paid or invalid state", nil)
		}

		payment := models.Payment{
			ID:      uuid.New(),
			OrderID: order.ID,
			Amount:  order.Total,
			Status:  models.PaymentPending,
			Method:  "mock_transfer",
		}

		if err := tx.Create(&payment).Error; err != nil {
			return utils.Internal("failed to create payment", err)
		}

		// mock success
		success := true

		if !success {
			payment.Status = models.PaymentFailed
			if err := tx.Save(&payment).Error; err != nil {
				fmt.Println("failed to update payment", err)
				return utils.Internal("failed to update payment", err)
			}
			return utils.BadRequest("payment failed", nil)
		}

		payment.Status = models.PaymentSuccess
		if err := tx.Save(&payment).Error; err != nil {
			return utils.Internal("failed to update payment", err)
		}

		order.Status = models.OrderPaid
		if err := tx.Save(&order).Error; err != nil {
			return utils.Internal("failed to update order", err)
		}

		return nil
	})

	if err != nil {
		c.Error(err)
		return
	}

	utils.Success(c, "payment success", nil)
}
