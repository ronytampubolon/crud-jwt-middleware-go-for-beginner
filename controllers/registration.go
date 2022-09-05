package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rony-tampubolon/sample-rest-api/models"
	"github.com/rony-tampubolon/sample-rest-api/utils"
)

type RegistrationInput struct {
	Name                 string `json:"name" binding:"required"`
	Email                string `json:"email" binding:"required"`
	Password             string `json:"password" binding:"required"`
	PasswordConfirmation string `json:"passwordConfirmation" binding:"required"`
}

func SignUp(c *gin.Context) {
	var input RegistrationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"err": err.Error(),
		})
		return
	}
	if input.Password != input.PasswordConfirmation {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"err": "password is not matched",
		})
		return
	}

	hash, _ := utils.HashPassword(input.Password)
	newUser := models.User{Name: input.Name, Email: input.Email, Password: hash, IsVerified: false}
	models.DB.Create(&newUser)

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "successfully registered",
		"data":    nil,
	})

}
