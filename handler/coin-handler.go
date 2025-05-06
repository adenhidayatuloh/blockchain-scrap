package handler

import (
	"blockchain-scrap/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type CoinHandler struct {
	Service service.CoinService
}

func NewCoinHandler(service service.CoinService) *CoinHandler {
	return &CoinHandler{Service: service}
}

func (h *CoinHandler) GetCoinDetail(c *gin.Context) {
	id := c.Param("id")
	contract := c.Param("contract-address")
	timeSkipStr := c.DefaultQuery("time-skip", "5m")
	const timeThreshold = 5 * time.Minute

	timeSkip, err := time.ParseDuration(timeSkipStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid time_skip format. Use formats like 30s, 5m, 1h."})
		return
	}

	if timeSkip < timeThreshold {
		c.JSON(http.StatusBadRequest, gin.H{"error": "time skip should more than 4 minutes"})
		return
	}

	result, err := h.Service.GetCoinDetail(id, contract, timeSkip)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *CoinHandler) GetAllCoins(c *gin.Context) {
	result, err := h.Service.GetAllCoins()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch data"})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *CoinHandler) StreamCoins(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Flush()

	ticker := time.NewTicker(7 * time.Second)
	defer ticker.Stop()

	sendCoinData := func() {
		coins, err := h.Service.GetAllCoins()
		if err != nil {
			return
		}
		c.SSEvent("message", coins)
	}

	sendCoinData()
	for {
		select {
		case <-c.Request.Context().Done():
			return
		case <-ticker.C:
			sendCoinData()
		}
	}
}
