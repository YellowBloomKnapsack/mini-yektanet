package handlers

import (
	"YellowBloomKnapsack/mini-yektanet/adserver/logic"

	"fmt"
	"net/http"
	"os"

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
	eventServerURL := "http://" + hostName + ":" + eventServerPort

	clickReqPath := os.Getenv("CLICK_REQ_PATH")
	impressionReqPath := os.Getenv("IMPRESSION_REQ_PATH")

	publisherUsername := c.Param("publisherUsername")

	clickToken, _ := logic.GenerateToken(logic.ClickType, chosenAd.ID, publisherUsername, chosenAd.Website)
	impressionToken, _ := logic.GenerateToken(logic.ImpressionType, chosenAd.ID, publisherUsername, chosenAd.Website)

	fmt.Println(clickToken)
	fmt.Println(impressionToken)

	c.JSON(http.StatusOK, gin.H{
		"image_link":       chosenAd.ImagePath,
		"title":            chosenAd.Text,
		"impression_link":  eventServerURL + "/" + impressionReqPath,
		"click_link":       eventServerURL + "/" + clickReqPath,
		"impression_token": impressionToken,
		"click_token":      clickToken,
	})
}
