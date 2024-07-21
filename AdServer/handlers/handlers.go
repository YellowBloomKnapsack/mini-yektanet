package handlers

import (
	"YellowBloomKnapsack/mini-yektanet/AdServer/logic"

	"net/http"
	"os"

	"github.com/beevik/guid"
	"github.com/gin-gonic/gin"
)

func GetAd(c *gin.Context) {
	chosenAd, err := logic.GetBestAd()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}

	eventServerPort := os.Getenv("EVENT_SERVER_PORT")
	hostName := os.Getenv("HOSTNAME")
	eventServerURL := hostName + ":" + eventServerPort

	clickReqPath := os.Getenv("CLICK_REQ_PATH")
	impressionReqPath := os.Getenv("IMPRESSION_REQ_PATH")

	clickToken := guid.NewString()
	impressionToken := guid.NewString()

	c.JSON(http.StatusOK, gin.H{
		"ad_id":            chosenAd.ID,
		"image_link":       chosenAd.ImagePath,
		"title":            chosenAd.Text,
		"impression_link":  eventServerURL + "/" + impressionReqPath,
		"click_link":       eventServerURL + "/" + clickReqPath,
		"impression_token": impressionToken,
		"click_token":      clickToken,
	})
}
