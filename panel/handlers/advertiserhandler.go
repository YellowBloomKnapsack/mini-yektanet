package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"

	"YellowBloomKnapsack/mini-yektanet/common/database"
	"YellowBloomKnapsack/mini-yektanet/common/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AdvertiserPanel(c *gin.Context) {
	// TODO: Get advertiser ID from session
	advertiserUserName := c.Param("username")

	var advertiser models.Advertiser
	result := database.DB.Preload("Ads").Where("username = ?", advertiserUserName).First(&advertiser)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			fmt.Println("No advertiser found with username %s, creating a new one.", advertiserUserName)
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

	c.HTML(http.StatusOK, "advertiser_panel.html", gin.H{
		"Balance":  advertiser.Balance,
		"Ads":      advertiser.Ads,
		"Username": advertiserUserName,
	})
}

func AddFunds(c *gin.Context) {
	advertiserUserName := c.Param("username")

	amount, _ := strconv.ParseInt(c.PostForm("amount"), 10, 64)

	database.DB.Model(&models.Advertiser{}).Where("username = ?", advertiserUserName).Update("balance", gorm.Expr("balance + ?", amount))

	c.Redirect(http.StatusSeeOther, fmt.Sprintf("/advertiser/%s/panel", advertiserUserName))
}

func CreateAd(c *gin.Context) {
	// TODO: Get advertiser ID from session
	advertiserUserName := c.Param("username")
	var advertiser models.Advertiser
	result := database.DB.Where("username = ?", advertiserUserName).First(&advertiser)
	if result.Error != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	title := c.PostForm("title")
	website := c.PostForm("website")
	bid, _ := strconv.ParseInt(c.PostForm("bid"), 10, 64)

	// Handle file upload
	file, _ := c.FormFile("image")
	// TODO: Save file and get path

	// Create a unique filename
	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)

	// Define the path where the image will be saved
	uploadDir := "static/uploads/"
	filepath := path.Join(uploadDir, filename)

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

	ad := models.Ad{
		Text:         title,
		ImagePath:    "/" + filepath, // Store the path relative to the server root
		Bid:          bid,
		AdvertiserID: advertiser.ID,
		Website: website,
	}

	database.DB.Create(&ad)

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
	adID, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var ad models.Ad
	database.DB.Preload("AdsInteractions").Preload("Advertiser").First(&ad, adID)

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
	})
}
