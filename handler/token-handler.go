package handler

import (
	"blockchain-scrap/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TokenHandler struct {
	service service.TokenService
}

func NewTokenHandler(s service.TokenService) *TokenHandler {
	return &TokenHandler{s}
}

// GetAllTokens godoc
// @Summary Get all tokens
// @Description Get paginated list of tokens with optional search
// @Tags token
// @Accept json
// @Produce json
// @Param limit query int false "Number of items per page (default: 10)"
// @Param page query int false "Page number (default: 1)"
// @Param search query string false "Search term"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /api/v1/tokens [get]
func (h *TokenHandler) GetAllTokens(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	pageStr := c.DefaultQuery("page", "1")
	search := c.Query("search")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	offset := (page - 1) * limit

	tokens, errService := h.service.GetAllTokens(limit, offset, search)
	if err != nil {

		c.JSON(errService.StatusCode(), gin.H{
			"message": errService.Message(),
			"error":   errService.Error(),
		})
		return

	}

	c.JSON(http.StatusOK, gin.H{
		"data": tokens,
		"pagination": gin.H{
			"total": tokens.Total,
			"page":  page,
			"limit": limit,
			"pages": (tokens.Total + int64(limit) - 1) / int64(limit),
		},
	})
}

// GetAccountInfo godoc
// @Summary Get account information
// @Description Get token account information by address
// @Tags token
// @Accept json
// @Produce json
// @Param address query string true "Account address"
// @Success 200 {object} dto.TokenAccountsResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/tokens/accounts/ [get]
func (h *TokenHandler) GetAccountInfo(c *gin.Context) {
	address := c.Query("address")

	if address == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Address is required",
		})
		return
	}

	result, errService := h.service.FetchAccountInfo(address)
	if errService != nil {

		c.JSON(errService.StatusCode(), gin.H{
			"error": errService.Error(),
		})
		return

	}

	c.JSON(http.StatusOK, result)
}
