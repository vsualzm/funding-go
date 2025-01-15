package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/vsualzm/funding-go/auth"
	"github.com/vsualzm/funding-go/campaign"
	"github.com/vsualzm/funding-go/handler"
	"github.com/vsualzm/funding-go/helper"
	"github.com/vsualzm/funding-go/user"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// var DB *gorm.DB

func main() {

	// cheking connection database
	// var err error

	dsn := "host=localhost user=postgres password=1234 dbname=startup_db port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	userRepository := user.NewRepository(db)

	// check package campaigns
	campaignRepository := campaign.NewRepository(db)

	campaigns, err := campaignRepository.FindByUserID(1)

	if err != nil {
		panic(err)
	}

	for _, campaign := range campaigns {

		fmt.Println(campaign.Name)
		fmt.Println(campaign.UserID)
		fmt.Println(campaign.Slug)
	}

	userService := user.NewService(userRepository)
	authService := auth.NewService()
	userHandler := handler.NewUserHandler(userService, authService)

	// api for gin

	router := gin.Default()

	api := router.Group("api/v1")

	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sesion", userHandler.Login)
	api.POST("/email_check", userHandler.CheckEmailAvailability)
	api.POST("/avatars", authMiddleware(authService, userService), userHandler.UploadAvatar)

	// nge checking save avatar
	router.Run(":8080")

}

func authMiddleware(authService auth.Service, userService user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if !strings.Contains(authHeader, "Bearer") {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		// bearer token
		tokenString := ""

		arrayToken := strings.Split(authHeader, " ")
		if len(arrayToken) == 2 {
			tokenString = arrayToken[1]
		}

		token, err := authService.ValidateToken(tokenString)

		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return

		}

		claim, ok := token.Claims.(jwt.MapClaims)

		if !ok || !token.Valid {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		userID := int(claim["user_id"].(float64))

		user, err := userService.GetUserByID(userID)

		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		c.Set("currentUser", user)

	}

}

// ambil nilai dari header Auth : Bearer token
// dari header auth kita ambil tokennya saja
// validasi token
// kita ambil user id
//  ambil user dari database
// set context isinya user
