package services

import (
	"project-e-commerce/models"

	"gorm.io/gorm"
	"time"
)

func GetTotalOrders(db *gorm.DB) (int64, error) {
	var count int64
	err := db.Model(&models.Order{}).Count(&count).Error
	return count, err
}
func GetTotalRevenue(db *gorm.DB) (float64, error) {
	var total float64

	err := db.Model(&models.Order{}).
		Select("COALESCE(SUM(total),0)").
		Where("status = ?", models.OrderCompleted).
		Scan(&total).Error

	return total, err
}

func GetOrdersByStatus(db *gorm.DB) (map[string]int64, error) {

	results := make(map[string]int64)

	rows, err := db.Model(&models.Order{}).
		Select("status, COUNT(*) as count").
		Group("status").
		Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var status string
		var count int64
		rows.Scan(&status, &count)
		results[status] = count
	}

	return results, nil
}

func GetMonthlyRevenue(db *gorm.DB) ([]map[string]interface{}, error) {

	rows, err := db.Raw(`
		SELECT 
			DATE_TRUNC('month', created_at) as month,
			SUM(total) as revenue
		FROM orders
		WHERE status = ?
		GROUP BY month
		ORDER BY month ASC
	`, models.OrderCompleted).Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []map[string]interface{}

	for rows.Next() {
		var month time.Time
		var revenue float64
		rows.Scan(&month, &revenue)

		results = append(results, map[string]interface{}{
			"month":   month,
			"revenue": revenue,
		})
	}

	return results, nil
}

func GetBestSellingProducts(db *gorm.DB) ([]map[string]interface{}, error) {

	rows, err := db.Raw(`
		SELECT 
			p.id,
			p.name,
			SUM(oi.quantity) as total_sold
		FROM order_items oi
		JOIN products p ON p.id = oi.product_id
		JOIN orders o ON o.id = oi.order_id
		WHERE o.status = ?
		GROUP BY p.id
		ORDER BY total_sold DESC
		LIMIT 5
	`, models.OrderCompleted).Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []map[string]interface{}

	for rows.Next() {
		var id string
		var name string
		var totalSold int64
		rows.Scan(&id, &name, &totalSold)

		results = append(results, map[string]interface{}{
			"id":         id,
			"name":       name,
			"total_sold": totalSold,
		})
	}

	return results, nil
}