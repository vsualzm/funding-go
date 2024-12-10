package main

import (
	"github.com/gin-gonic/gin"
	"github.com/vsualzm/funding-go/auth"
	"github.com/vsualzm/funding-go/handler"
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
	userService := user.NewService(userRepository)
	authService := auth.NewService()

	userHandler := handler.NewUserHandler(userService, authService)

	// api for gin

	router := gin.Default()

	api := router.Group("api/v1")

	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sesion", userHandler.Login)
	api.POST("/email_check", userHandler.CheckEmailAvailability)
	api.POST("/avatars", userHandler.UploadAvatar)

	// nge checking save avatar
	router.Run(":8080")

}

// func Handler(c *gin.Context) {

// 	dsn := "host=localhost user=postgres password=1234 dbname=startup_db port=5432 sslmode=disable TimeZone=Asia/Shanghai"
// 	DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

// 	if err != nil {
// 		panic("failed to connect database")
// 	}

// 	// get data from database
// 	var users []user.User
// 	DB.Find(&users)
// 	c.JSON(200, users)

// }

// input
// handler menangkap inputan dari user mapping ke strcut
// service mapping ke strcut
// repository menyimpan ke db
// db
