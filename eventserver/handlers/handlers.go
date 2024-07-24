package handlers

import (
	"os"
	"time"
	"context"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"encoding/base64"
	"net/http"
	"github.com/redis/go-redis/v9"

	"YellowBloomKnapsack/mini-yektanet/common/tokenhandler"
	"YellowBloomKnapsack/mini-yektanet/eventserver/worker"
)

var age int
var redisClient redis.Client

type EventServerHandler struct {
	addToken func(string, bool)
	isTokenPresent func(string, bool) bool
	tokenHandler tokenhandler.TokenHandlerInterface
	workerService worker.WorkerInterface
}

func NewEventServerHandler(tokenHandler tokenhandler.TokenHandlerInterface, workerService worker.WorkerInterface) *EventServerHandler {
	workerService.Start()
	age, _ = strconv.Atoi(os.Getenv("REDIS_AGE_HOURS"))
	redisClient = *redis.NewClient(&redis.Options{
		Addr:	  os.Getenv("REDIS_URL"),
        Password: "", // no password set
        DB:		  0,  // use default DB
	})
	return &EventServerHandler{
		addToken: addToken,
		isTokenPresent: isTokenPresent,
		tokenHandler: tokenHandler,
		workerService: workerService,
	}
}

type TokenRequest struct {
	Token string `json:"token"`
}

func (h *EventServerHandler) PostClick(c *gin.Context) {
	privateKey := os.Getenv("PRIVATE_KEY")
	key, _ := base64.StdEncoding.DecodeString(privateKey)

	var req TokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token := req.Token
	data, err := h.tokenHandler.VerifyToken(token, key)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	present := h.isTokenPresent(token, true)
	if !present {
		h.addToken(token, true)
		h.workerService.InvokeClickEvent(data, time.Now())
	}

	c.Redirect(http.StatusMovedPermanently, data.RedirectPath)
}

func (h *EventServerHandler) PostImpression(c *gin.Context) {
	privateKey := os.Getenv("PRIVATE_KEY")
	key, _ := base64.StdEncoding.DecodeString(privateKey)

	var req TokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token := req.Token
	data, err := h.tokenHandler.VerifyToken(token, key)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	present := h.isTokenPresent(token, false)
	if !present {
		h.addToken(token, false)
		h.workerService.InvokeImpressionEvent(data, time.Now())
	}
}

func addToken(token string, isClick bool) {
	ctx := context.Background()
	var err error
	if isClick {
		err = redisClient.Set(ctx, "click:"+token, true, time.Duration(age)*time.Hour).Err()
	} else {
		err = redisClient.Set(ctx, "impression:"+token, true, time.Duration(age)*time.Hour).Err()
	}
	if err != nil {
		log.Print(err)
	}
}

func isTokenPresent(token string, isClick bool) bool {
	ctx := context.Background()
	var key string
	if isClick {
		key = "click:"+token
	} else {
		key = "impression:"+token
	}
	present, err := redisClient.Exists(ctx, key).Result()
	if err != nil {
		log.Print(err)
		return false
	}
	if present == 1 {
		return true
	} else {
		return false
	}
}