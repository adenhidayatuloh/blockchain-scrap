package handler

import (
	"blockchain-scrap/entity"
	"blockchain-scrap/pkg/errs"
	"blockchain-scrap/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// BlockchainHandler handles blockchain data requests
type BlockchainHandler struct {
	blockchainSvc service.BlockchainService
}

// NewBlockchainHandler creates a new instance of BlockchainHandler
func NewBlockchainHandler(svc service.BlockchainService) *BlockchainHandler {
	return &BlockchainHandler{blockchainSvc: svc}
}

// GetBlockchainDetailByContractAddress gets blockchain details by contract address
// GetBlockchainDetailByContractAddress godoc
// @Summary Get blockchain details by contract address
// @Description Get detailed information about a blockchain token using its contract address
// @Tags blockchain
// @Accept json
// @Produce json
// @Param contract-address path string true "Contract Address"
// @Param time-skip query string false "Time interval for data (default: 5m)" default(5m)
// @Success 200 {object} dto.ContractAddressResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /coins/v2/{contract-address} [get]
func (h *BlockchainHandler) GetBlockchainDetailByContractAddress(c *gin.Context) {
	contractAddress := c.Param("contract-address")
	timeSkipStr := c.DefaultQuery("time-skip", "5m")
	const minTimeSkip = 5 * time.Minute
	//userData := c.MustGet("userData").(*entity.User)

	timeSkip, err := time.ParseDuration(timeSkipStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid time_skip format. Use formats like 30s, 5m, 1h."})
		return
	}

	if timeSkip < minTimeSkip {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Time interval must be more than 4 minutes"})
		return
	}

	result, errService := h.blockchainSvc.GetBlockchainDetailByContractAddress(c.Request.Context(), contractAddress, timeSkip)
	if errService != nil {

		c.JSON(errService.StatusCode(), gin.H{"error": errService.Message()})

		return
	}
	c.JSON(http.StatusOK, result)
}

// GetBlockchainDetailByIDAndContractAddress gets blockchain details by ID and contract address
// GetBlockchainDetailByIDAndContractAddress godoc
// @Summary Get blockchain details by ID and contract address
// @Description Get detailed information about a blockchain token using blockchain ID and contract address
// @Tags blockchain
// @Accept json
// @Produce json
// @Param blockchain-id path string true "Blockchain ID"
// @Param contract-address path string true "Contract Address"
// @Param time-skip query string false "Time interval for data (default: 5m)" default(5m)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /coins/{blockchain-id}/{contract-address} [get]
func (h *BlockchainHandler) GetBlockchainDetailByIDAndContractAddress(c *gin.Context) {
	contractAddress := c.Param("contract-address")
	blockchainID := c.Param("blockchain-id")
	timeSkipStr := c.DefaultQuery("time-skip", "5m")
	const minTimeSkip = 5 * time.Minute

	timeSkip, err := time.ParseDuration(timeSkipStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid time_skip format. Use formats like 30s, 5m, 1h."})
		return
	}

	if timeSkip < minTimeSkip {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Time interval must be more than 4 minutes"})
		return
	}

	result, err := h.blockchainSvc.GetBlockchainDetailByContractAddressAndID(blockchainID, contractAddress, timeSkip)
	if err != nil {
		if messageErr, ok := err.(errs.MessageErr); ok {
			c.JSON(messageErr.StatusCode(), gin.H{"error": messageErr.Message()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error occurred"})
		}
		return
	}
	c.JSON(http.StatusOK, result)
}

// GetAllBlockchains gets all blockchain data
// GetAllBlockchains godoc
// @Summary Get all blockchains
// @Description Get a list of all available blockchains
// @Tags blockchain
// @Accept json
// @Produce json
// @Success 200 {array} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /api/v1/blockchains [get]
func (h *BlockchainHandler) GetAllBlockchains(c *gin.Context) {
	result, err := h.blockchainSvc.GetAllBlockchains()
	if err != nil {
		c.JSON(err.StatusCode(), gin.H{"error": err.Message()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// GetAllBlockchainSearchesByUserID gets all search history by user ID
// GetAllBlockchainSearchesByUserID godoc
// @Summary Get all blockchain searches by user ID
// @Description Get search history for a specific user
// @Tags blockchain
// @Accept json
// @Produce json
// @Success 200 {array} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /api/v1/searches [get]
func (h *BlockchainHandler) GetAllBlockchainSearchesByUserID(c *gin.Context) {
	userData := c.MustGet("userData").(*entity.User)
	result, err := h.blockchainSvc.FindByUserID(c.Request.Context(), userData.ID)
	if err != nil {
		c.JSON(err.StatusCode(), gin.H{"error": err.Message()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// GetBlockchainSearchByID gets search history by ID
// GetBlockchainSearchByID godoc
// @Summary Get blockchain search by ID
// @Description Get specific search history by search ID
// @Tags blockchain
// @Accept json
// @Produce json
// @Param search-id path string true "Search ID (UUID)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/searches/{search-id} [get]
func (h *BlockchainHandler) GetBlockchainSearchByID(c *gin.Context) {
	searchID := c.Param("search-id")

	searchUUID, err := uuid.Parse(searchID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID must be a UUID"})
		return
	}

	result, err := h.blockchainSvc.FindByID(c.Request.Context(), searchUUID)
	if err != nil {
		if messageErr, ok := err.(errs.MessageErr); ok {
			c.JSON(messageErr.StatusCode(), gin.H{"error": messageErr.Message()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error occurred"})
		}
		return
	}
	c.JSON(http.StatusOK, result)
}

// StreamBlockchains streams blockchain data in real-time using Server-Sent Events
// StreamBlockchains godoc
// @Summary Stream blockchain data
// @Description Stream real-time blockchain data using Server-Sent Events
// @Tags blockchain
// @Accept json
// @Produce text/event-stream
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /api/v1/blockchains/stream [get]
func (h *BlockchainHandler) StreamBlockchains(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Flush()

	ticker := time.NewTicker(7 * time.Second)
	defer ticker.Stop()

	sendBlockchainData := func() {
		blockchains, err := h.blockchainSvc.GetAllBlockchains()
		if err != nil {
			c.SSEvent("error", gin.H{"error": err.Message()})
			return
		}
		c.SSEvent("message", blockchains)
	}

	sendBlockchainData()
	for {
		select {
		case <-c.Request.Context().Done():
			return
		case <-ticker.C:
			sendBlockchainData()
		}
	}
}
