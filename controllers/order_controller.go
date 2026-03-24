package controllers

import (
	"errors"
	"project-e-commerce/config"
	"project-e-commerce/models"
	"project-e-commerce/services"
	"project-e-commerce/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UpdateOrderStatusRequest struct {
	Status string `json:"status" binding:"required"`
}
// Checkout godoc
// @Summary Checkout order
// @Tags orders
// @Security BearerAuth
// @Produce json
// @Success 200 {object} models.Order
// @Router /orders/checkout [post]
func Checkout(c *gin.Context) {

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

	err = config.DB.Transaction(func(tx *gorm.DB) error {

		var cart models.Cart
		if err := tx.Preload("CartItems.Product").
			Where("user_id = ?", userID).
			First(&cart).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return utils.BadRequest("cart not found", nil)
			}
			return utils.Internal("database error", err)
		}

		if len(cart.CartItems) == 0 {
			return utils.BadRequest("cart is empty", nil)
		}

		var total float64 = 0

		for _, item := range cart.CartItems {
			if item.Product.Stock < item.Quantity {
				return utils.BadRequest("insufficient stock", nil)
			}
			total += item.Product.Price * float64(item.Quantity)
		}

		order := models.Order{
			ID:     uuid.New(),
			UserID: userID,
			Total:  total,
			Status: "pending",
		}

		if err := tx.Create(&order).Error; err != nil {
			return err
		}

		for _, item := range cart.CartItems {

			orderItem := models.OrderItem{
				ID:        uuid.New(),
				OrderID:   order.ID,
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
				Price:     item.Product.Price,
			}

			if err := tx.Create(&orderItem).Error; err != nil {
				return err
			}
			if err := tx.Model(&models.Product{}).
				Where("id = ?", item.ProductID).
				Update("stock", gorm.Expr("stock - ?", item.Quantity)).Error; err != nil {
				return err
			}
		}

		// 🔥 ลบ cart items หลังสร้าง order
		if err := tx.Where("cart_id = ?", cart.ID).
			Delete(&models.CartItem{}).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		c.Error(utils.BadRequest(err.Error(), err))
		return
	}

	utils.Success(c, "checkout success", nil)
}
// GetMyOrders godoc
// @Summary Get my orders
// @Tags orders
// @Security BearerAuth
// @Produce json
// @Success 200 {array} models.Order
// @Router /orders [get]
func GetMyOrders(c *gin.Context) {

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

	var orders []models.Order


	if err := config.DB.
		Preload("OrderItems").
		Where("user_id = ?", userID).
    	Order("created_at DESC").
    	Find(&orders).Error; err != nil {

		c.Error(utils.Internal("database error", err))
		return
	}

	utils.Success(c, "my orders", orders)
}
// GetOrderDetail godoc
// @Summary Get order detail
// @Tags orders
// @Security BearerAuth
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} models.Order
// @Router /orders/{id} [get]
func GetOrderDetail(c *gin.Context) {

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

	orderID := c.Param("id")

	var order models.Order

	if err := config.DB.
		Preload("OrderItems").
		Where("id = ? AND user_id = ?", orderID, userID).
		First(&order).Error; err != nil {

		c.Error(utils.BadRequest("order not found", err))
		return
	}

	utils.Success(c, "order detail", order)
}
// UpdateOrderStatus godoc
// @Summary Update order status
// @Tags orders
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Param body body UpdateOrderStatusRequest true "Update order status data"
// @Success 200 {object} map[string]interface{}
// @Router /orders/{id}/status [patch]
func UpdateOrderStatus(c *gin.Context) {
	orderID := c.Param("id")


	var input UpdateOrderStatusRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.Error(utils.BadRequest("invalid input", err))
		return
	}

	// ดึง user จาก auth middleware
	currentUserID := c.MustGet("user_id").(uuid.UUID)

	// 🔥 เริ่ม transaction
	err := config.DB.Transaction(func(tx *gorm.DB) error {

	var order models.Order
	if err := tx.First(&order, "id = ?", orderID).Error; err != nil {
		return utils.BadRequest("order not found", nil)
	}

	return services.ChangeOrderStatus(tx, &order, input.Status, currentUserID)
})

	if err != nil {
		c.Error(err)
		return
	}

	utils.Success(c, "order updated", nil)
}
// CancelOrder godoc
// @Summary Cancel order
// @Tags orders
// @Security BearerAuth
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} map[string]interface{}
// @Router /orders/{id}/cancel [put]
func CancelOrder(c *gin.Context) {

	// ✅ ดึง user id แบบ safe
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

	orderID := c.Param("id")

	err = config.DB.Transaction(func(tx *gorm.DB) error {

		var order models.Order

		// ✅ หา order ของ user คนนี้เท่านั้น
		if err := tx.Preload("OrderItems").
			Where("id = ? AND user_id = ?", orderID, userID).
			First(&order).Error; err != nil {

			if errors.Is(err, gorm.ErrRecordNotFound) {
				return utils.BadRequest("order not found", nil)
			}

			return utils.Internal("database error", err)
		}

		// ✅ เช็คว่า cancel ได้ไหม
		if !services.CanTransition(order.Status, models.OrderCancelled) {
			return utils.BadRequest("cannot cancel this order", nil)
		}

		// ✅ คืน stock แบบ atomic
		for _, item := range order.OrderItems {

			if err := tx.Model(&models.Product{}).
				Where("id = ?", item.ProductID).
				Update("stock", gorm.Expr("stock + ?", item.Quantity)).Error; err != nil {

				return utils.Internal("failed to restore stock", err)
			}
		}

		// ✅ update status
		order.Status = models.OrderCancelled

		if err := tx.Save(&order).Error; err != nil {
			return utils.Internal("failed to update order", err)
		}

		return nil
	})

	if err != nil {
		c.Error(err)
		return
	}

	utils.Success(c, "order cancelled successfully", nil)
}
// GetAllOrders godoc
// @Summary Get all orders
// @Tags orders
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page"
// @Param limit query int false "Limit"
// @Success 200 {array} models.Order
// @Router /admin/orders [get]
func GetAllOrders(c *gin.Context) {

	var orders []models.Order
	page := 1
	limit := 20
	offset := (page - 1) * limit

	if err := config.DB.
		Preload("OrderItems").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&orders).Error; err != nil {

		c.Error(utils.Internal("database error", err))
		return
	}

	utils.Success(c, "all orders", orders)
}

