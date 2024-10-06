package main

import (
	"donation/controllers"
	"donation/helper"
	"donation/repository"
	"donation/services"
	"log"
	"net/http"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:@tcp(127.0.0.1:3306)/donation?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err.Error())
	}

	userRepository := repository.NewUserRepository(db)
	donationRepository := repository.NewDonationRepository(db)
	transactionRepository := repository.NewTransactionRepository(db)

	userService := services.NewUserService(userRepository)
	donationService := services.NewDonationService(donationRepository)
	paymentService := services.NewPaymentServicce(transactionRepository, donationRepository)
	transactionService := services.NewTransactionService(transactionRepository, donationRepository, paymentService)
	authService := services.NewJWTService()

	userController := controllers.NewUserController(userService, authService)
	donationController := controllers.NewDonationController(donationService)
	transactionController := controllers.NewTransactionController(transactionService, paymentService)

	router := gin.New()

	config := cors.DefaultConfig()
	config.AllowHeaders = append(config.AllowHeaders, "Authorization")
	config.AllowAllOrigins = true
	router.Use(cors.New(config))

	router.Static("/images", "./images")
	api := router.Group("/api")
	api.POST("/register", userController.RegisterUser)
	api.POST("/register_admin", userController.RegisterAdmin)
	api.POST("/login", userController.Login)
	api.POST("/email_check", userController.CheckEmailAvailability)
	api.POST("/avatar", authMiddleware(authService, userService), userController.UploadAvatar)
	api.POST("/update_profile", authMiddleware(authService, userService), userController.UpdateProfile)
	api.POST("/update_password", authMiddleware(authService, userService), userController.UpdatePassword)
	api.GET("/donations", donationController.GetDonations)
	api.GET("/donations/:id", donationController.GetDonation)
	api.GET("/donations/user/:user_id", authMiddleware(authService, userService), donationController.GetDonationsByUserID)
	api.POST("/donations", authMiddleware(authService, userService), donationController.CreateDonation)
	api.PUT("/donations/:id", authMiddleware(authService, userService), donationController.UpdateDonation)
	api.DELETE("/donations/:id", authMiddleware(authService, userService), donationController.DeleteDonation)
	api.GET("/donations/:id/transactions", authMiddleware(authService, userService), transactionController.GetDonationTransations)
	api.GET("/transactions", authMiddleware(authService, userService), transactionController.GetUserTransactions)
	api.GET("/transactions/:id", transactionController.GetTransactionByID)
	api.POST("/transactions", authMiddleware(authService, userService), transactionController.CreateTransaction)
	api.POST("/transactions/notification", transactionController.GetNotification)

	router.Run()
}

func authMiddleware(authService services.AuthService, userService services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if !strings.Contains(authHeader, "Bearer") {
			response := helper.ResponseJSON("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		tokenString := ""
		arrayToken := strings.Split(authHeader, " ")
		if len(arrayToken) == 2 {
			tokenString = arrayToken[1]
		} else {
			response := helper.ResponseJSON("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		token, err := authService.ValidateToken(tokenString)
		if err != nil {
			response := helper.ResponseJSON("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		claim, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			response := helper.ResponseJSON("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		userID := int(claim["user_id"].(float64))

		user, err := userService.GetUserByID(userID)
		if err != nil {
			response := helper.ResponseJSON("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		c.Set("currentUser", user)
	}
}
