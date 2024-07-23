package handlers

import (
	"YellowBloomKnapsack/mini-yektanet/common/models"
	"YellowBloomKnapsack/mini-yektanet/panel/database"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func PublisherPanel(c *gin.Context) {
	username := c.Param("username")

	var publisher models.Publisher
	result := database.DB.Where("username = ?", username).Preload("AdsInteraction").Preload("AdsInteraction.Ad").First(&publisher)
	if result.Error != nil {
		fmt.Println("Info:", result.Error)
		publisher = models.Publisher{
			Username: username,
			Balance:  0,
		}
		createResult := database.DB.Create(&publisher)
		if createResult.Error != nil {
			fmt.Println("Error creating Publisher:", createResult.Error)
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		fmt.Println("New Publisher created:", publisher)
		// c.JSON(http.StatusNotFound, gin.H{"error": "Publisher not found"})
		// return
	}

	script := fmt.Sprintf("<script src='%s:%s/js/script.js'></script>",
		os.Getenv("HOSTNAME"),
		os.Getenv("PUBLISHER_WEBSITE_PORT"))

	yektanetPortionString := os.Getenv("YEKTANET_PORTION")

	// Convert the value to an integer
	yektanetPortion, err := strconv.Atoi(yektanetPortionString)
	if err != nil || yektanetPortion < 0 || yektanetPortion > 100 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error parsing YEKTANET_PORTION environment variable: %v\n", err)})
		return
	}

	// Prepare data for the chart
	chartData := prepareChartData(publisher.AdsInteraction, yektanetPortion)

	c.HTML(http.StatusOK, "publisher_panel.html", gin.H{
		"Publisher": publisher,
		"Script":    script,
		"ChartData": chartData,
	})
}
func WithdrawPublisherBalance(c *gin.Context) {
	username := c.Param("username")

	var publisher models.Publisher
	result := database.DB.Where("username = ?", username).First(&publisher)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Publisher not found"})
		return
	}

	amountStr := c.PostForm("amount")
	amount, err := strconv.ParseInt(amountStr, INTBASE, INTBIT64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid amount"})
		return
	}

	if amount <= 0 || amount > publisher.Balance {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid withdrawal amount"})
		return
	}

	// amount = min(amount, publisher.Balance)
	publisher.Balance -= amount
	database.DB.Save(&publisher)

	c.JSON(http.StatusOK, gin.H{
		"message":    fmt.Sprintf("Withdrawn amount: %d", amount),
		"newBalance": publisher.Balance,
	})
}

func prepareChartData(interactions []models.AdsInteraction, yektanetPortion int) map[string]interface{} {
	// Group interactions by day
	dailyData := make(map[time.Time]struct {
		Impressions int
		Clicks      int
		Revenue     int64
	})

	for _, interaction := range interactions {
		day := interaction.ClickTime.Truncate(24 * time.Hour)
		data := dailyData[day]
		if interaction.Type == int(models.Impression) {
			data.Impressions++
		} else if interaction.Type == int(models.Click) {
			data.Clicks++
			data.Revenue += interaction.Ad.Bid
		}
		data.Revenue = data.Revenue * int64(100-yektanetPortion) / 100
		dailyData[day] = data
	}

	// Convert to slices for JSON serialization
	var dates []string
	var impressions []int
	var clicks []int
	var revenues []int64

	for date := time.Now().AddDate(0, 0, -30); date.Before(time.Now()); date = date.AddDate(0, 0, 1) {
		data := dailyData[date.Truncate(24*time.Hour)]
		dates = append(dates, date.Format("2006-01-02"))
		impressions = append(impressions, data.Impressions)
		clicks = append(clicks, data.Clicks)
		revenues = append(revenues, data.Revenue)
	}

	return map[string]interface{}{
		"Dates":       dates,
		"Impressions": impressions,
		"Clicks":      clicks,
		"Revenues":    revenues,
	}
}
