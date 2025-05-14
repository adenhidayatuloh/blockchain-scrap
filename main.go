package main

import (
	"blockchain-scrap/handler"
	"blockchain-scrap/infra"
	"blockchain-scrap/repository"
	"blockchain-scrap/service"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Initialize database
	db := infra.GetDBInstance()
	if err := infra.AutoMigrate(db); err != nil {
		return
	}

	// Initialize router
	router := gin.Default()

	// Configure CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error load env file: %s", err)
	}

	apiPort := os.Getenv("APP_PORT")

	// Initialize repositories
	blockchainSearchRepo := repository.NewBlockchainSearchRepository(db)
	tokenRepo := repository.NewTokenRepository(db)
	userRepo := repository.NewUserRepository(db)

	// Initialize services
	blockchainService := service.NewBlockchainService(blockchainSearchRepo)
	tokenService := service.NewTokenService(tokenRepo)
	userService := service.NewUserService(userRepo)
	swapService := service.NewSwapService(tokenRepo, tokenService)

	// Initialize handlers
	blockchainHandler := handler.NewBlockchainHandler(blockchainService)
	tokenHandler := handler.NewTokenHandler(tokenService)
	userHandler := handler.NewUserHandler(userService)
	swapHandler := handler.NewSwapHandler(swapService)

	//Default first api
	router.GET("/coins/v2/:contract-address", blockchainHandler.GetBlockchainDetailByContractAddress)
	router.GET("/coins/:blockchain-id/:contract-address", blockchainHandler.GetBlockchainDetailByIDAndContractAddress)

	// API v1 routes
	v1 := router.Group("/api/v1")

	{
		// Public routes
		// Auth routes
		auth := v1.Group("/auth")
		{
			auth.POST("/register", userHandler.Register)
			auth.POST("/login", userHandler.Login)
		}

		// Protected routes
		protected := v1.Group("")
		//protected.Use(userService.Authentication())
		{
			// Token routes
			tokens := protected.Group("/tokens")
			{
				tokens.GET("", tokenHandler.GetAllTokens)
				tokens.GET("/accounts/", tokenHandler.GetAccountInfo)
			}

			// Swap routes
			swaps := protected.Group("/swaps")
			{
				swaps.POST("", swapHandler.Swap)
				swaps.POST("/submit", swapHandler.Submit)
				swaps.POST("/quote", swapHandler.GetCurrencySwap)
			}

			// Blockchain routes
			blockchains := protected.Group("/blockchains")
			{
				blockchains.GET("", blockchainHandler.GetAllBlockchains)
				blockchains.GET("/stream", blockchainHandler.StreamBlockchains)
				blockchains.GET("/:contract-address", blockchainHandler.GetBlockchainDetailByContractAddress)
			}

			// Search history routes
			searches := protected.Group("/searches")
			{
				searches.GET("", blockchainHandler.GetAllBlockchainSearchesByUserID)
				searches.GET("/:search-id", blockchainHandler.GetBlockchainSearchByID)
			}
		}
	}

	router.Run(":" + apiPort)
}
