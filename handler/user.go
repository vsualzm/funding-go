package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vsualzm/funding-go/helper"
	"github.com/vsualzm/funding-go/user"
)

type userHandler struct {
	userService user.Service
}

func NewUserHandler(userService user.Service) *userHandler {
	return &userHandler{userService}
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

	formatter := user.FormatUser(newUser, "tokennn")

	response := helper.APIResponse("Account has been registered", http.StatusOK, "success", formatter)

	// response handler
	c.JSON(http.StatusOK, response)

}
