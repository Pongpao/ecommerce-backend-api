package controllers

import (
	"errors"
	"project-e-commerce/config"
	"project-e-commerce/models"
	"project-e-commerce/services"
	"project-e-commerce/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)
// CreateReview godoc
// @Summary Create product review
// @Description User can review a purchased product
// @Tags reviews
// @Accept json
// @Produce json
// @Param review body models.Review true "Review data"
// @Success 200 {object} models.Review
// @Failure 400 {object} utils.ErrorResponse
// @Failure 403 {object} utils.ErrorResponse
// @Router /reviews [post]

func CreateReview(c *gin.Context) {

	userIDStr := c.MustGet("user_id").(string)

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.Error(utils.Internal("invalid user id", err))
		return
	}
	var body struct {
		ProductID string `json:"product_id" binding:"required"`
		Rating    int    `json:"rating" binding:"required,min=1,max=5"`
		Comment   string `json:"comment"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.Error(utils.BadRequest("invalid input request", err))
		return
	}

	productUUID, err := uuid.Parse(body.ProductID)
	if err != nil {
		c.Error(utils.BadRequest("invalid product id", err))
		return
	}

	// check product exists
	var product models.Product
	if err := config.DB.First(&product, "id = ?", productUUID).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Error(utils.NotFound("product not found", nil))
			return
		}

		c.Error(utils.Internal("database error", err))
		return
	}

	// check purchased
	purchased, err := services.HasPurchasedProduct(
		config.DB,
		userID.String(),
		body.ProductID,
	)
	if err != nil {
		c.Error(utils.Internal("check failed", err))
		return
	}

	if !purchased {
		c.Error(utils.Forbidden("you must purchase this product before reviewing", nil))
		return
	}

	// check already reviewed
	var existing models.Review
	err = config.DB.
		Where("user_id = ? AND product_id = ?", userID, productUUID).
		First(&existing).Error

	if err == nil {
		c.Error(utils.BadRequest("already reviewed", nil))
		return
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		c.Error(utils.Internal("database error", err))
		return
	}

	review := models.Review{
		UserID:    userID,
		ProductID: productUUID,
		Rating:    body.Rating,
		Comment:   body.Comment,
	}

	err = config.DB.Transaction(func(tx *gorm.DB) error {

		if err := tx.Create(&review).Error; err != nil {
			if strings.Contains(err.Error(), "duplicate") {
				return utils.BadRequest("already reviewed", err)
			}
			return err
		}

		var avg float64

		if err := tx.
			Model(&models.Review{}).
			Where("product_id = ?", productUUID).
			Select("AVG(rating)").
			Scan(&avg).Error; err != nil {
			return err
		}

		if err := tx.
			Model(&models.Product{}).
			Where("id = ?", productUUID).
			Update("average_rating", avg).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		c.Error(err)
		return
	}
	utils.Success(c, "review created", review)
}
