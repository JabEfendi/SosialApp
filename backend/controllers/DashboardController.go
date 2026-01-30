package controllers

import (
	"net/http"
	"backend/db"

	"github.com/gin-gonic/gin"
)

func CorporateIncomeChart(c *gin.Context) {
	corporateID := c.MustGet("corporate_id").(uint)

	type Result struct {
		Month  string `json:"month"`
		Total  int64  `json:"total"`
	}

	var results []Result

	db.DB.Raw(`
		SELECT 
			DATE_FORMAT(created_at, '%Y-%m') AS month,
			SUM(amount) AS total
		FROM token_ledgers
		WHERE user_id = ?
		  AND source_type = 'room_revenue_release'
		GROUP BY month
		ORDER BY month ASC
	`, corporateID).Scan(&results)

	c.JSON(http.StatusOK, gin.H{
		"data": results,
	})
}

func CorporateWithdrawRequestChart(c *gin.Context) {
	corporateID := c.MustGet("corporate_id").(uint)

	type Result struct {
		Status string `json:"status"`
		Total  int64  `json:"total"`
	}

	var results []Result

	db.DB.Raw(`
		SELECT 
			status,
			COUNT(*) AS total
		FROM withdraw_requests
		WHERE corporate_id = ?
		GROUP BY status
	`, corporateID).Scan(&results)

	c.JSON(http.StatusOK, gin.H{
		"data": results,
	})
}

func UserWithdrawAccumulationChart(c *gin.Context) {

	type Result struct {
		Month string `json:"month"`
		Total int64  `json:"total"`
	}

	var results []Result

	db.DB.Raw(`
		SELECT 
			DATE_FORMAT(created_at, '%Y-%m') AS month,
			SUM(amount) AS total
		FROM withdraw_requests
		WHERE requested_by = 'user'
		  AND status = 'approved'
		GROUP BY month
		ORDER BY month ASC
	`).Scan(&results)

	c.JSON(http.StatusOK, gin.H{
		"data": results,
	})
}

func CorporateDashboard(c *gin.Context) {
	corporateID := c.MustGet("corporate_id").(uint)

	var income []map[string]interface{}
	var wdCorp []map[string]interface{}
	var wdUser []map[string]interface{}

	db.DB.Raw(`
		SELECT DATE_FORMAT(created_at, '%Y-%m') month, SUM(amount) total
		FROM token_ledgers
		WHERE user_id = ? AND source_type = 'room_revenue_release'
		GROUP BY month
	`, corporateID).Scan(&income)

	db.DB.Raw(`
		SELECT status, COUNT(*) total
		FROM withdraw_requests
		WHERE corporate_id = ?
		GROUP BY status
	`, corporateID).Scan(&wdCorp)

	db.DB.Raw(`
		SELECT DATE_FORMAT(created_at, '%Y-%m') month, SUM(amount) total
		FROM withdraw_requests
		WHERE requested_by = 'user' AND status = 'approved'
		GROUP BY month
	`).Scan(&wdUser)

	c.JSON(http.StatusOK, gin.H{
		"income_corporate": income,
		"wd_corporate":     wdCorp,
		"wd_user":          wdUser,
	})
}