package handlers

import (
	"encoding/base64"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"

	"YellowBloomKnapsack/mini-yektanet/adserver/grafana"
	"YellowBloomKnapsack/mini-yektanet/adserver/logic"
	"YellowBloomKnapsack/mini-yektanet/common/models"
	"YellowBloomKnapsack/mini-yektanet/common/tokenhandler"
)

type AdServerHandler struct {
	logicService logic.LogicInterface
	tokenHandler tokenhandler.TokenHandlerInterface
}

func NewAdServerHandler(logicService logic.LogicInterface, tokenHandler tokenhandler.TokenHandlerInterface) *AdServerHandler {
	logicService.StartTicker()

	return &AdServerHandler{
		logicService: logicService,
		tokenHandler: tokenHandler,
	}
}

func (h *AdServerHandler) GetAd(c *gin.Context) {
	grafana.AdRequestTotal.Inc()
	timer := prometheus.NewTimer(grafana.AdRequestDuration)
	defer timer.ObserveDuration()

	chosenAd, err := h.logicService.GetBestAd()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}

	eventServerPort := os.Getenv("EVENT_SERVER_PORT")
	hostName := os.Getenv("EVENT_SERVER_PUBLIC_HOSTNAME")
	eventServerURL := "http://" + hostName + ":" + eventServerPort

	clickReqPath := os.Getenv("CLICK_REQ_PATH")
	impressionReqPath := os.Getenv("IMPRESSION_REQ_PATH")

	publisherId, _ := strconv.Atoi(c.Param("publisherId"))

	privateKey := os.Getenv("PRIVATE_KEY")
	key, _ := base64.StdEncoding.DecodeString(privateKey)

	clickToken, clickErr := h.tokenHandler.GenerateToken(models.Click, chosenAd.ID, uint(publisherId), chosenAd.Bid, chosenAd.Website, key)
	if clickErr != nil {
		grafana.TokenGenerationErrorsTotal.WithLabelValues("click").Inc()
	}

	impressionToken, impressionErr := h.tokenHandler.GenerateToken(models.Impression, chosenAd.ID, uint(publisherId), chosenAd.Bid, chosenAd.Website, key)

	if impressionErr != nil {
		grafana.TokenGenerationErrorsTotal.WithLabelValues("impression").Inc()
	}

	c.JSON(http.StatusOK, gin.H{
		"image_link":       chosenAd.ImagePath,
		"title":            chosenAd.Text,
		"impression_link":  eventServerURL + "/" + impressionReqPath,
		"click_link":       eventServerURL + "/" + clickReqPath,
		"impression_token": impressionToken,
		"click_token":      clickToken,
	})
	grafana.AdServedTotal.Inc()
}

func (h *AdServerHandler) BrakeAd(c *gin.Context) {
	adId, err := strconv.Atoi(c.Param("adId"))
	if err != nil {
		grafana.InvalidAdBrakeRequestsTotal.Inc()
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid adId"})
		return
	}

	h.logicService.BrakeAd(uint(adId))
}
