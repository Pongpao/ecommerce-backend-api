package services

import (
	"project-e-commerce/models"

	"gorm.io/gorm"
)

type ProductFilter struct {
	Search     string
	MinPrice   float64
	MaxPrice   float64
	MinRating  float64
	InStock    *bool
	CategoryID string
	Sort       string
	Page       int
	Limit      int
}

func GetFilteredProducts(db *gorm.DB, filter ProductFilter) ([]models.Product, int64, error) {

	var products []models.Product
	var total int64

	query := db.Model(&models.Product{})

	if filter.Search != "" {
		query = query.Where("name ILIKE ?", "%"+filter.Search+"%")
	}

	if filter.MinPrice > 0 {
		query = query.Where("price >= ?", filter.MinPrice)
	}

	if filter.MaxPrice > 0 {
		query = query.Where("price <= ?", filter.MaxPrice)
	}

	if filter.MinRating > 0 {
		query = query.Where("average_rating >= ?", filter.MinRating)
	}

	if filter.InStock != nil && *filter.InStock {
		query = query.Where("stock > 0")
	}

	if filter.CategoryID != "" {
		query = query.Where("category_id = ?", filter.CategoryID)
	}

	countQuery := query.Session(&gorm.Session{})
	countQuery.Count(&total)

	switch filter.Sort {
	case "price_asc":
		query = query.Order("price ASC")
	case "price_desc":
		query = query.Order("price DESC")
	case "rating":
		query = query.Order("average_rating DESC")
	default:
		query = query.Order("created_at DESC")
	}

	offset := (filter.Page - 1) * filter.Limit

	err := query.
		Offset(offset).
		Limit(filter.Limit).
		Find(&products).Error

	return products, total, err
}