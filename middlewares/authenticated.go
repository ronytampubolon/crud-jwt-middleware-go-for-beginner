package middlewares

import (
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/rony-tampubolon/sample-rest-api/entities"
	"github.com/rony-tampubolon/sample-rest-api/utils"
)

func AuthenticatedUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		println("Enter here")
		authorizationHeader := c.GetHeader("Authorization")
		if !strings.Contains(authorizationHeader, "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"err": "not authorized",
			})
			return
		}

		tokenString := strings.Replace(authorizationHeader, "Bearer ", "", -1)

		claims := &entities.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return utils.SecretKey, nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"status":  http.StatusUnauthorized,
					"message": "not authorized",
					"data":    nil,
				})
				return
			}
			if !token.Valid {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"status":  http.StatusBadRequest,
					"message": "Invalid Token",
					"data":    nil,
				})
				return
			}
		}
		c.Next()

	}
}
