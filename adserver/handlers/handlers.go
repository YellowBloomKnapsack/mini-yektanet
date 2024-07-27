package handlers

import (
	"encoding/base64"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"YellowBloomKnapsack/mini-yektanet/adserver/logic"
	"YellowBloomKnapsack/mini-yektanet/common/dto"
	"YellowBloomKnapsack/mini-yektanet/common/tokenhandler"
)

type AdServerHandler struct {
	logicService  logic.LogicInterface
	tokenHandler  tokenhandler.TokenHandlerInterface
	brakeDuration time.Duration
}

func NewAdServerHandler(logicService logic.LogicInterface, tokenHandler tokenhandler.TokenHandlerInterface) *AdServerHandler {
	logicService.StartTicker()
	brakeDuration, _ := strconv.Atoi(os.Getenv("BRAKE_DURATION_SECS"))

	return &AdServerHandler{
		logicService:  logicService,
		tokenHandler:  tokenHandler,
		brakeDuration: time.Duration(brakeDuration) * time.Second,
	}
}

func (h *AdServerHandler) GetAd(c *gin.Context) {
	chosenAd, err := h.logicService.GetBestAd()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}

	eventServerPort := os.Getenv("EVENT_SERVER_PORT")
	hostName := os.Getenv("EVENT_SERVER_HOSTNAME")
	eventServerURL := "http://" + hostName + ":" + eventServerPort

	clickReqPath := os.Getenv("CLICK_REQ_PATH")
	impressionReqPath := os.Getenv("IMPRESSION_REQ_PATH")

	publisherUsername := c.Param("publisherUsername")

	privateKey := os.Getenv("PRIVATE_KEY")
	key, _ := base64.StdEncoding.DecodeString(privateKey)

	clickToken, _ := h.tokenHandler.GenerateToken(dto.ClickType, chosenAd.ID, publisherUsername, chosenAd.Website, key)
	impressionToken, _ := h.tokenHandler.GenerateToken(dto.ImpressionType, chosenAd.ID, publisherUsername, chosenAd.Website, key)

	c.JSON(http.StatusOK, gin.H{
		"image_link":       chosenAd.ImagePath,
		"title":            chosenAd.Text,
		"impression_link":  eventServerURL + "/" + impressionReqPath,
		"click_link":       eventServerURL + "/" + clickReqPath,
		"impression_token": impressionToken,
		"click_token":      clickToken,
	})
}

func (h *AdServerHandler) BrakeAd(c *gin.Context) {
	adId, err := strconv.Atoi(c.Param("adId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid adId"})
		return
	}

	h.logicService.BrakeAd(uint(adId), h.brakeDuration)
}
