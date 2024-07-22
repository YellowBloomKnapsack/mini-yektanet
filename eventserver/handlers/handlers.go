package handlers

import (
    "fmt"
    "strconv"
    "net/http"

    "github.com/gin-gonic/gin"
)

// "math/rand"
// "net/http"
// "time"

// "YellowBloomKnapsack/mini-yektanet/common/database"
// "YellowBloomKnapsack/mini-yektanet/common/models"


var clickTokens = make(map[string]bool, 0)
var impressionTokens = make(map[string]bool, 0)

type EventPayload struct {
    AdID string `json:"ad_id"`
    PublisherID string `json:"publisher_id"`
    Token string `json:"token"`
    RedirectPath string `json:"redirect_path"`
}

func PostClick(c *gin.Context) {
    var requestBody EventPayload 
    if err := c.BindJSON(&requestBody); err != nil {
        fmt.Println("Something wrong happened")
        return
    }
    
    adID, _ := strconv.Atoi(requestBody.AdID)
    publisherID, _ := strconv.Atoi(requestBody.PublisherID)

    _, present := clickTokens[requestBody.Token]
    if !present {
        clickTokens[requestBody.Token] = true
        fmt.Println("jadid")
        // 1) request panel api
    }

    // redirect anyways
	c.Redirect(http.StatusMovedPermanently, requestBody.RedirectPath)

    fmt.Println(adID)
    fmt.Println(publisherID)
}
