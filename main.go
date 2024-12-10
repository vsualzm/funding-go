package main

import (
	"github.com/gin-gonic/gin"
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

	userHandler := handler.NewUserHandler(userService)

	// api for gin

	router := gin.Default()

	api := router.Group("api/v1")

	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sesion", userHandler.Login)
	api.POST("/email_check", userHandler.CheckEmailAvailability)
	api.POST("/avatars", userHandler.UploadAvatar)

	// nge checking save avatar
	router.Run(":8080")

	// userInput := user.RegisterUserInput{}
	// userInput.Name = "test simpan dari service"
	// userInput.Occupation = "informa"
	// userInput.Email = "tX6oJ@example.com"
	// userInput.Password = "colokan password"

	// userService.RegisterUser(userInput)

	// user := user.User{
	// 	Name:         "Visual",
	// 	Occupation:   "Programmer",
	// 	Email:        "q4qK1@example.com",
	// 	PasswordHash: "password",
	// }

	// cheking save function
	// userRepository.Save(user)

	// fmt.Println("Connection to Database Successful")

	// var user []user.User

	// fmt.Println(len(user))

	// DB.Find(&user)

	// fmt.Println(len(user))

	// for _, u := range user {
	// 	fmt.Println(u.Name)
	// 	fmt.Println(u.Email)
	// 	fmt.Println(u.PasswordHash)
	// 	fmt.Println("================")
	// }

	// router := gin.Default()
	// // router.GET("/users", Handler)

	// router.Run(":8000")

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
