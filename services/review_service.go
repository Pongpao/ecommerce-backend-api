package services

import (
	"project-e-commerce/models"

	"gorm.io/gorm"
)

func HasPurchasedProduct(db *gorm.DB, userID, productID string) (bool, error) {
	var count int64

	err := db.
		Model(&models.OrderItem{}).
		Joins("JOIN orders ON orders.id = order_items.order_id").
		Where("orders.user_id = ? AND order_items.product_id = ? AND orders.status = ?",
			userID, productID, models.OrderCompleted).
		Count(&count).Error

	return count > 0, err
}
func CalculateAverageRating(db *gorm.DB, productID string) (float64, error) {
	var avg float64

	err := db.Model(&models.Review{}).
		Select("AVG(rating)").
		Where("product_id = ?", productID).
		Scan(&avg).Error

	return avg, err
}