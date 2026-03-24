package controllers

import (
	"errors"
	"project-e-commerce/config"
	"project-e-commerce/models"
	"project-e-commerce/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)
type AddToCartRequest struct {
	ProductID string `json:"product_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required"`
}
type UpdateCartItemRequest struct {
	ProductID string `json:"product_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required"`
}

// AddToCart godoc
// @Summary Add product to cart
// @Tags cart
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body AddToCartRequest true "Add to cart data"
// @Success 200 {object} map[string]interface{}
// @Router /cart/items [post]
func AddToCart(c *gin.Context) {

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

	var input AddToCartRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.Error(utils.BadRequest("invalid input", err))
		return
	}

	productUUID, err := uuid.Parse(input.ProductID)
	if err != nil {
		c.Error(utils.Internal("invalid product id", err))
		return
	}
	if input.Quantity <= 0 {
		c.Error(utils.BadRequest("quantity must be greater than 0", nil))
		return
	}

	err = config.DB.Transaction(func(tx *gorm.DB) error {

		var product models.Product

		if err := tx.First(&product, "id = ?", productUUID).Error; err != nil {
			return err
		}

		if product.Stock < input.Quantity {
			return utils.BadRequest("not enough stock", nil)
		}
		// 🔎 หา cart
		var cart models.Cart

		err := tx.Where("user_id = ?", userID).
			First(&cart).Error

		if errors.Is(err, gorm.ErrRecordNotFound) {

			cart = models.Cart{
				ID:     uuid.New(),
				UserID: userID,
			}

			if err := tx.Create(&cart).Error; err != nil {
				return err
			}

		} else if err != nil {
			return err
		}
		// 🔎 เพิ่ม cart item
		var cartItem models.CartItem

		err = tx.Where("cart_id = ? AND product_id = ?", cart.ID, product.ID).
			First(&cartItem).Error

		if errors.Is(err, gorm.ErrRecordNotFound) {

			cartItem = models.CartItem{
				ID:        uuid.New(),
				CartID:    cart.ID,
				ProductID: product.ID,
				Quantity:  input.Quantity,
			}

			return tx.Create(&cartItem).Error

		}

		if err != nil {
			return err
		}

		newQty := cartItem.Quantity + input.Quantity

		if newQty > product.Stock {
			return utils.BadRequest("not enough stock", nil)
		}

		cartItem.Quantity = newQty

		return tx.Model(&cartItem).Update("quantity", cartItem.Quantity).Error
	})

	if err != nil {
		c.Error(utils.BadRequest(err.Error(), err))
		return
	}

	utils.Success(c, "added to cart", nil)
}
// GetCart godoc
// @Summary Get cart
// @Tags cart
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /cart [get]
func GetCart(c *gin.Context) {

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

	var cart models.Cart
	err = config.DB.Where("user_id = ?", userID).
		First(&cart).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Success(c, "cart is empty", gin.H{
				"items": []interface{}{},
			})
			return
		}
		c.Error(utils.Internal("database error", err))
		return
	}

	var items []models.CartItem
	config.DB.
		Preload("Product").
		Where("cart_id = ?", cart.ID).
		Find(&items)

	utils.Success(c, "cart items", items)
}

// UpdateCartItem godoc
// @Summary Update cart item
// @Tags cart
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body UpdateCartItemRequest true "Update cart item data"
// @Success 200 {object} map[string]interface{}
// @Router /cart/items [put]
func UpdateCartItem(c *gin.Context) {

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

	var input struct {
		ProductID string `json:"product_id" binding:"required"`
		Quantity  int    `json:"quantity" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.Error(utils.BadRequest("invalid request body", err))
		return
	}

	if input.Quantity < 0 {
		c.Error(utils.BadRequest("quantity cannot be negative", nil))
		return
	}

	productUUID, err := uuid.Parse(input.ProductID)
	if err != nil {
		c.Error(utils.Internal("invalid product id", err))
		return
	}

	err = config.DB.Transaction(func(tx *gorm.DB) error {

		// 1️⃣ หา cart
		var cart models.Cart
		if err := tx.Where("user_id = ?", userID).
			First(&cart).Error; err != nil {
			return utils.BadRequest("cart not found", err)
		}

		// 2️⃣ หา cart item
		var cartItem models.CartItem
		if err := tx.Where("cart_id = ? AND product_id = ?",
			cart.ID, productUUID).
			First(&cartItem).Error; err != nil {
			return utils.BadRequest("item not in cart", err)
		}

		// 🔥 ถ้า quantity = 0 → ลบทิ้ง
		if input.Quantity == 0 {
			return tx.Delete(&cartItem).Error
		}

		// 3️⃣ หา product
		var product models.Product
		if err := tx.First(&product,
			"id = ?", productUUID).Error; err != nil {
			return err
		}

		// 4️⃣ เช็ค stock
		if input.Quantity > product.Stock {
			return utils.BadRequest("not enough stock", err)
		}

		// 5️⃣ update quantity
		cartItem.Quantity = input.Quantity
		return tx.Save(&cartItem).Error

	})

	if err != nil {
		c.Error(utils.BadRequest(err.Error(), err))
		return
	}

	utils.Success(c, "cart updated", nil)
}
