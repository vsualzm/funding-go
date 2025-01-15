package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vsualzm/funding-go/auth"
	"github.com/vsualzm/funding-go/helper"
	"github.com/vsualzm/funding-go/user"
)

type userHandler struct {
	userService user.Service
	authService auth.Service
}

func NewUserHandler(userService user.Service, authService auth.Service) *userHandler {
	return &userHandler{userService, authService}
}

func (h *userHandler) RegisterUser(c *gin.Context) {
	// tangkap data dari input user
	// map input dari user ke struct registeruserinput
	// struct di atas kita parsing sebagai parameter service

	// dapat input dari user
	var input user.RegisterUserInput
	err := c.ShouldBindJSON(&input)
	if err != nil {

		// var errors []string

		// // mapping looping error handling request bair rapih tidak satu line kalau seperti ini
		// for _, e := range err.(validator.ValidationErrors) {
		// 	errors = append(errors, e.Error())
		// }
		// errorMessage := gin.H{"errors": errors}

		errors := helper.FormatValidationError(err)
		response := helper.APIResponse("Register Account Failed", http.StatusUnprocessableEntity, "success", errors)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// panggil service registeruser
	newUser, err := h.userService.RegisterUser(input)
	if err != nil {
		response := helper.APIResponse("Register Account Failed", http.StatusBadRequest, "success", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// generate token
	token, err := h.authService.GenerateToken(newUser.ID)
	if err != nil {
		response := helper.APIResponse("Register Account Failed", http.StatusBadRequest, "success", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := user.FormatUser(newUser, token)

	response := helper.APIResponse("Account has been registered", http.StatusOK, "success", formatter)

	// response handler
	c.JSON(http.StatusOK, response)

}

func (h *userHandler) Login(c *gin.Context) {

	// user input Email
	// input di tangkap handler
	// mapping dari input user ke input struct
	// input struct nya passing ke service
	// di service mencari dengan repository user dengan email x
	// mencocokan password'

	var input user.LoginInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIResponse("Login Failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	loggedinUser, err := h.userService.Login(input)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}
		response := helper.APIResponse("Login Failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	tokenLogin, err := h.authService.GenerateToken(loggedinUser.ID)

	if err != nil {
		response := helper.APIResponse("Login Failed", http.StatusBadRequest, "success", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := user.FormatUser(loggedinUser, tokenLogin)
	response := helper.APIResponse("Successfuly loggedin", http.StatusOK, "success", formatter)
	c.JSON(http.StatusOK, response)

}

func (h *userHandler) CheckEmailAvailability(c *gin.Context) {

	// user input Email
	var input user.CheckEmailInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIResponse("Email Checking Failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// check email dengan service email
	isEmailAvaible, err := h.userService.CheckEmailAvailability(input)

	if err != nil {
		errorMessage := gin.H{"errors": "Server Error"}
		response := helper.APIResponse("Email Checking Failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	data := gin.H{
		"is_available": isEmailAvaible,
	}

	metaMessage := "Email has been registered"

	if isEmailAvaible {
		metaMessage = "Email is available"
	}

	response := helper.APIResponse(metaMessage, http.StatusOK, "success", data)

	c.JSON(http.StatusOK, response)

}

func (h *userHandler) UploadAvatar(c *gin.Context) {

	// ambil data dari form file
	file, err := c.FormFile("avatar")
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to upload avatar", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// dapatkan current user id
	currentUser := c.MustGet("currentUser").(user.User)
	userID := currentUser.ID

	// simpan ke folder images dengan format nama seperti ini
	path := fmt.Sprintf("images/%d-%s", userID, file.Filename)
	err = c.SaveUploadedFile(file, path)

	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to upload avatar", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// simpan ke db
	_, err = h.userService.SaveAvatar(userID, path)

	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to upload avatar", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// berikan response ini dalam data
	data := gin.H{
		"is_uploaded": true,
	}

	response := helper.APIResponse("Avatar successfully uploaded", http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)

}
