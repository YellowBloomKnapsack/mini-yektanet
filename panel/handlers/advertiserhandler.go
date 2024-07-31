package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"YellowBloomKnapsack/mini-yektanet/common/models"
	"YellowBloomKnapsack/mini-yektanet/panel/database"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	INTBASE      = 10
	INTBIT32     = 32
	INTBIT64     = 64
	itemsPerPage = 47
)

func AdvertiserPanel(c *gin.Context) {
	advertiserUserName := c.Param("username")

	var advertiser models.Advertiser
	result := database.DB.Preload("Ads").Preload("Ads.Keywords").Where("username = ?", advertiserUserName).First(&advertiser)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			fmt.Printf("No advertiser found with username %s, creating a new one.\n", advertiserUserName)
			newAdvertiser := models.Advertiser{
				Username: advertiserUserName,
				Balance:  0,
			}
			createResult := database.DB.Create(&newAdvertiser)
			if createResult.Error != nil {
				fmt.Println("Error creating advertiser:", createResult.Error)
				c.AbortWithStatus(http.StatusInternalServerError)
			}
			fmt.Println("New advertiser created:", newAdvertiser)

		} else {
			fmt.Println("Error:", result.Error)
			c.AbortWithStatus(http.StatusInternalServerError)
		}
	}

	// Pagination logic
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	var totalTransactions int64
	database.DB.Model(&models.Transaction{}).Where("customer_id = ? AND customer_type = ?", advertiser.ID, models.Customer_Advertiser).Count(&totalTransactions)
	totalPages := int((totalTransactions + itemsPerPage - 1) / itemsPerPage)

	var transactions []models.Transaction
	database.DB.Model(&models.Transaction{}).Where("customer_id = ? AND customer_type = ?", advertiser.ID, models.Customer_Advertiser).
		Offset((page - 1) * itemsPerPage).
		Limit(itemsPerPage).
		Find(&transactions)

	c.HTML(http.StatusOK, "advertiser_panel.html", gin.H{
		"Balance":      advertiser.Balance,
		"Ads":          advertiser.Ads,
		"Username":     advertiserUserName,
		"Transactions": transactions,
		"TotalPages":   totalPages,
		"CurrentPage":  page + 1,
	})
}

func AddFunds(c *gin.Context) {
	advertiserUserName := c.Param("username")

	amount, err := strconv.ParseInt(c.PostForm("amount"), INTBASE, INTBIT64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid amount"})
		return
	}

	if amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid amount"})
		return
	}
	tx := database.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process deposit"})
		}
	}()

	// Fetch the advertiser
	var advertiser models.Advertiser
	if err := tx.Where("username = ?", advertiserUserName).First(&advertiser).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch advertiser"})
		return
	}
	//if advertiser.Balance < 0 && advertiser.Balance+amount >= 0 {
	//	go NotifyAdsUpdate()
	//}
	// Update the advertiser's balance
	advertiser.Balance += amount
	if err := tx.Save(&advertiser).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update advertiser balance"})
		return
	}

	// Create a new transaction record
	transaction := models.Transaction{
		CustomerID:   advertiser.ID,
		CustomerType: models.Customer_Advertiser,
		Amount:       amount,
		Income:       false,
		Successful:   true,
		Time:         time.Now(),
		Description:  "charge wallet",
	}
	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction"})
		return
	}

	tx.Commit()

	c.Redirect(http.StatusSeeOther, fmt.Sprintf("/advertiser/%s/panel", advertiserUserName))
}

func CreateAd(c *gin.Context) {
	advertiserUserName := c.Param("username")
	var advertiser models.Advertiser
	result := database.DB.Where("username = ?", advertiserUserName).First(&advertiser)
	if result.Error != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	title := c.PostForm("title")
	website := c.PostForm("website")
	bid, _ := strconv.ParseInt(c.PostForm("bid"), INTBASE, INTBIT64)
	if bid <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bid amount"})
		return
	}

	file, _ := c.FormFile("image")
	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)
	uploadDir := "static/uploads/"
	filepath := path.Join(uploadDir, filename)

	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
		return
	}

	if err := c.SaveUploadedFile(file, filepath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	ad := models.Ad{
		Text:         title,
		ImagePath:    "/" + filepath,
		Bid:          bid,
		AdvertiserID: advertiser.ID,
		Website:      website,
	}

	database.DB.Create(&ad)

	// Extract and save keywords
	keywords := []string{
		c.PostForm("keyword1"),
		c.PostForm("keyword2"),
		c.PostForm("keyword3"),
		c.PostForm("keyword4"),
	}

	keywordString := strings.Join(keywords, ",")

	keywordRecord := models.Keyword{
		AdID:     ad.ID,
		Keywords: keywordString,
	}
	database.DB.Create(&keywordRecord)

	c.Redirect(http.StatusSeeOther, fmt.Sprintf("/advertiser/%s/panel", advertiserUserName))
}

func ToggleAd(c *gin.Context) {
	advertiserUserName := c.Param("username")
	adID, _ := strconv.ParseUint(c.PostForm("ad_id"), 10, 32)

	var ad models.Ad
	database.DB.First(&ad, adID)

	ad.Active = !ad.Active
	database.DB.Save(&ad)

	c.Redirect(http.StatusSeeOther, fmt.Sprintf("/advertiser/%s/panel", advertiserUserName))
}

func AdReport(c *gin.Context) {
	adID, _ := strconv.ParseUint(c.Param("id"), INTBASE, INTBIT32)

	var ad models.Ad
	database.DB.Preload("AdsInteractions").Preload("Keywords").Preload("Advertiser").First(&ad, adID)

	impressions := 0
	clicks := 0
	for _, interaction := range ad.AdsInteractions {
		if interaction.Type == int(models.Impression) {
			impressions++
		} else if interaction.Type == int(models.Click) {
			clicks++
		}
	}

	ctr := float64(0)
	if impressions > 0 {
		ctr = float64(clicks) / float64(impressions) * 100
	}

	c.HTML(http.StatusOK, "ad_report.html", gin.H{
		"Ad":          ad,
		"Impressions": impressions,
		"Clicks":      clicks,
		"CTR":         ctr,
		"TotalCost":   ad.TotalCost,
		"Website":     ad.Website,
		"Keywords":    ad.Keywords,
	})
}

func HandleEditAd(c *gin.Context) {
	username := c.Param("username")
	adID, err := strconv.Atoi(c.PostForm("ad_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ad ID"})
		return
	}

	// Find the advertiser
	var advertiser models.Advertiser
	if err := database.DB.Where("username = ?", username).First(&advertiser).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Advertiser not found"})
		return
	}

	// Find the ad
	var ad models.Ad
	if err := database.DB.Where("id = ? AND advertiser_id = ?", adID, advertiser.ID).First(&ad).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ad not found"})
		return
	}

	// Update ad details
	ad.Text = c.PostForm("text")
	ad.Website = c.PostForm("website")
	bid, err := strconv.ParseInt(c.PostForm("bid"), INTBASE, INTBIT64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bid amount"})
		return
	}
	ad.Bid = bid

	// Handle image upload if a new image is provided
	file, _ := c.FormFile("image")
	if file != nil {
		// Create a unique filename
		filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)

		// Define the path where the image will be saved
		uploadDir := "static/uploads/"
		filepath := path.Join(uploadDir, filename)
		removeFileIfExists(ad.ImagePath)

		// Ensure the upload directory exists
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
			return
		}

		// Save the file
		if err := c.SaveUploadedFile(file, filepath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
			return
		}
		ad.ImagePath = "/" + filepath
	}

	// Save the updated ad
	if err := database.DB.Save(&ad).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update ad"})
		return
	}

	c.Redirect(http.StatusSeeOther, "/advertiser/"+username+"/panel")
}

func removeFileIfExists(filePath string) error {
	// Check if the file exists
	_, err := os.Stat(filePath)
	if err == nil {
		// File exists, so remove it
		err := os.Remove(filePath)
		if err != nil {
			return fmt.Errorf("failed to remove file: %v", err)
		}
		fmt.Printf("File %s has been removed\n", filePath)
	} else if os.IsNotExist(err) {
		// File doesn't exist, so no need to remove
		fmt.Printf("File %s does not exist\n", filePath)
	} else {
		// Some other error occurred
		return fmt.Errorf("error checking file: %v", err)
	}
	return nil
}
