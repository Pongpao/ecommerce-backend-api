package controllers

import (
	"errors"
	"project-e-commerce/config"
	"project-e-commerce/models"
	"project-e-commerce/services"
	"project-e-commerce/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)
type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required"`
	Stock       int     `json:"stock" binding:"required"`
}
type UpdateProductRequest struct {
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	Price       *float64 `json:"price"`
	Stock       *int     `json:"stock"`
}
// CreateProduct godoc
// @Summary Create product
// @Tags products
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body CreateProductRequest true "Create product data"
// @Success 200 {object} models.Product
// @Router /products [post]
func CreateProduct(c *gin.Context) {
	var input CreateProductRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.Error(utils.BadRequest("invalid input", err))
		return
	}

	if input.Price <= 0 {
		c.Error(utils.BadRequest("price must be greater than 0", nil))
		return
	}

	product := models.Product{
		ID:          uuid.New(),
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		Stock:       input.Stock,
	}

	if err := config.DB.Create(&product).Error; err != nil {
		c.Error(utils.Internal("failed to create product", err))
		return
	}

	utils.Success(c, "product created", product)
}

func GetProductByID(c *gin.Context) {
	productID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.Error(utils.BadRequest("invalid product id", err))
		return
	}

	var product models.Product
	if err := config.DB.First(&product, "id = ?", productID).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Error(utils.NotFound("product not found", nil))
			return
		}

		c.Error(utils.Internal("database error", err))
		return
	}

	utils.Success(c, "product found", product)
}
// UpdateProduct godoc
// @Summary Update product
// @Tags products
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param body body CreateProductRequest true "Update product data"
// @Success 200 {object} models.Product
// @Router /products/{id} [put]
func UpdateProduct(c *gin.Context) {
	productID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.Error(utils.BadRequest("invalid product id", err))
		return
	}

	var product models.Product
	if err := config.DB.First(&product, "id = ?", productID).Error; err != nil {
		c.Error(utils.NotFound("product not found", err))
		return
	}

	var input UpdateProductRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.Error(utils.BadRequest("invalid input", err))
		return
	}

	updates := map[string]interface{}{}

	if input.Name != nil {
		updates["name"] = *input.Name
	}
	if input.Description != nil {
		updates["description"] = *input.Description
	}
	if input.Price != nil {
		updates["price"] = *input.Price
	}
	if input.Stock != nil {
		updates["stock"] = *input.Stock
	}

	if err := config.DB.Model(&product).Updates(updates).Error; err != nil {
		c.Error(utils.Internal("failed to update product", err))
		return
	}
	utils.Success(c, "product updated", product)
}

// DeleteProduct godoc
// @Summary Delete product
// @Tags products
// @Security BearerAuth
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} map[string]interface{}
// @Router /products/{id} [delete]
func DeleteProduct(c *gin.Context) {
	productID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.Error(utils.BadRequest("invalid product id", err))
		return
	}

	var product models.Product
	if err := config.DB.First(&product, "id = ?", productID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Error(utils.NotFound("product not found", nil))
			return
		}
		c.Error(utils.Internal("database error", err))
		return
	}

	result := config.DB.Delete(&models.Product{}, "id = ?", productID)

	if result.Error != nil {
		c.Error(utils.Internal("failed to delete", result.Error))
		return
	}

	if result.RowsAffected == 0 {
		c.Error(utils.NotFound("product not found", nil))
		return
	}

	utils.Success(c, "product deleted", nil)
}
// GetProducts godoc
// @Summary Get products
// @Tags products
// @Produce json
// @Param page query int false "Page"
// @Param limit query int false "Limit"
// @Success 200 {array} models.Product
// @Router /products [get]
func GetProducts(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		c.Error(utils.BadRequest("invalid page", nil))
		return
	}
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		c.Error(utils.BadRequest("invalid limit", nil))
		return
	}

	var minPrice float64
	if minPriceStr := c.Query("min_price"); minPriceStr != "" {
		minPrice, err = strconv.ParseFloat(minPriceStr, 64)
		if err != nil {
			c.Error(utils.BadRequest("invalid min_price", nil))
			return
		}
	}
	var maxPrice float64
	if maxPriceStr := c.Query("max_price"); maxPriceStr != "" {
		maxPrice, err = strconv.ParseFloat(maxPriceStr, 64)
		if err != nil {
			c.Error(utils.BadRequest("invalid max_price", nil))
			return
		}
	}
	var minRating float64
	if minRatingStr := c.Query("min_rating"); minRatingStr != "" {
		minRating, err = strconv.ParseFloat(minRatingStr, 64)
		if err != nil {
			c.Error(utils.BadRequest("invalid min_rating", nil))
			return
		}
	}
	var inStock *bool
	if val := c.Query("in_stock"); val != "" {
		parsed, err := strconv.ParseBool(val)
		if err != nil {
			c.Error(utils.BadRequest("invalid in_stock", nil))
			return
		}
		inStock = &parsed
	}

	filter := services.ProductFilter{
		Search:     c.Query("search"),
		MinPrice:   minPrice,
		MaxPrice:   maxPrice,
		MinRating:  minRating,
		InStock:    inStock,
		CategoryID: c.Query("category_id"),
		Sort:       c.Query("sort"),
		Page:       page,
		Limit:      limit,
	}

	products, total, err := services.GetFilteredProducts(config.DB, filter)
	if err != nil {
		c.Error(utils.Internal("failed to get products", err))
		return
	}

	utils.Success(c, "products fetched", gin.H{
		"total":    total,
		"page":     page,
		"limit":    limit,
		"products": products,
	})
}
