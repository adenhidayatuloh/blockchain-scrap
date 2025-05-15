package handler

import (
	"blockchain-scrap/dto"
	"blockchain-scrap/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SwapHandler struct {
	Service service.SwapService
}

func NewSwapHandler(s service.SwapService) *SwapHandler {
	return &SwapHandler{Service: s}
}

// Swap godoc
// @Summary Create swap transaction
// @Description Create a new token swap transaction
// @Tags swap
// @Accept json
// @Produce json
// @Param request body dto.SwapRequest true "Swap Request"
// @Success 200 {object} dto.SwapResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/swaps [post]
func (h *SwapHandler) Swap(c *gin.Context) {
	var req dto.SwapRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transaction, err := h.Service.GetSwapTransaction(req)
	if err != nil {
		c.JSON(err.StatusCode(), gin.H{"error": err.Message()})
		return
	}

	c.JSON(http.StatusOK, dto.SwapResponse{Transaction: transaction})
}

// Swap godoc
// @Summary post quote transaction
// @Description post quote transaction
// @Tags swap
// @Accept json
// @Produce json
// @Param request body dto.SwapRequest true "Swap Request"
// @Success 200 {object} dto.SwapResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/swaps/quote [post]
func (h *SwapHandler) GetCurrencySwap(c *gin.Context) {
	var req dto.SwapRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transaction, err := h.Service.GetCurrencySwap(req)
	if err != nil {
		c.JSON(err.StatusCode(), gin.H{"error": err.Message()})
		return
	}

	c.JSON(http.StatusOK, transaction)
}

// Submit godoc
// @Summary Submit swap transaction
// @Description Submit a swap transaction for processing
// @Tags swap
// @Accept json
// @Produce json
// @Param request body dto.SubmitRequest true "Submit Request"
// @Success 200 {object} dto.SubmitResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/swaps/submit [post]
func (h *SwapHandler) Submit(c *gin.Context) {
	var req dto.SubmitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	signature, err := h.Service.SubmitTransaction(req)
	if err != nil {
		c.JSON(err.StatusCode(), gin.H{"error": err.Message()})
		return
	}

	c.JSON(http.StatusOK, dto.SubmitResponse{Signature: signature})
}
