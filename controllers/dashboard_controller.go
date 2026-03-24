package controllers

import (
	"project-e-commerce/config"
	"project-e-commerce/services"
	"project-e-commerce/utils"

	"github.com/gin-gonic/gin"
)

func GetDashboard(c *gin.Context) {

	totalOrders, err := services.GetTotalOrders(config.DB)
	if err != nil {
		c.Error(err)
		return
	}

	totalRevenue, err := services.GetTotalRevenue(config.DB)
	if err != nil {
		c.Error(err)
		return
	}

	ordersByStatus, err := services.GetOrdersByStatus(config.DB)
	if err != nil {
		c.Error(err)
		return
	}

	monthlyRevenue, err := services.GetMonthlyRevenue(config.DB)
	if err != nil {
		c.Error(err)
		return
	}

	bestProducts, err := services.GetBestSellingProducts(config.DB)
	if err != nil {
		c.Error(err)
		return
	}

	utils.Success(c, "dashboard fetched", gin.H{
		"total_orders":     totalOrders,
		"total_revenue":    totalRevenue,
		"orders_by_status": ordersByStatus,
		"monthly_revenue":  monthlyRevenue,
		"best_products":    bestProducts,
	})
}