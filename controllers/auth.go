package controllers

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/rony-tampubolon/sample-rest-api/entities"
	"github.com/rony-tampubolon/sample-rest-api/models"
	"github.com/rony-tampubolon/sample-rest-api/utils"
)

type CredentialInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func SignIn(c *gin.Context) {
	var input CredentialInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"err": err.Error(),
		})
		return
	}
	var selecteduser models.User

	// check if user is exist
	if err := models.DB.Where("email = ?", input.Email).First(&selecteduser).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}
	// check password
	if match := utils.CheckPasswordHash(input.Password, selecteduser.Password); match == false {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "wrong password",
		})
		return
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &entities.Claims{
		Username: input.Email,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(utils.SecretKey)
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": err.Error(),
		})
		return
	}
	// set Cookie Token
	http.SetCookie(c.Writer, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})
}

func RefreshToken(c *gin.Context) {

	claims := &entities.Claims{}

	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Bad Request",
			"data":    "Token still valid",
		})
		return
	}

	// Now, create a new token for the current use, with a renewed expiration time
	expirationTime := time.Now().Add(5 * time.Minute)
	claims.ExpiresAt = expirationTime.Unix()
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, error := newToken.SignedString(utils.SecretKey)
	if error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Internal Server Error",
			"data":    error.Error(),
		})
		return
	}

	// Set the new token as the users `token` cookie
	http.SetCookie(c.Writer, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})
}
