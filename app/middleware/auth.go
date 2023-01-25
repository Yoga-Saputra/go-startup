package middleware

import (
	"bwastartup/app/auth"
	"bwastartup/app/helper"
	"bwastartup/app/users"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(authService auth.Service, userService users.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.Contains(authHeader, "Bearer") {
			responseFailedMiddleware(c)
			return
		}

		// Bearer token
		tokenString := ""
		arrayToken := strings.Split(authHeader, " ")
		if len(arrayToken) == 2 {
			tokenString = arrayToken[1]
		}

		token, err := authService.ValidateToken(tokenString)

		if err != nil {
			responseFailedMiddleware(c)
			return
		}

		claim, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			responseFailedMiddleware(c)
			return
		}

		userId := int(claim["user_id"].(float64))
		user, err := userService.GetUserByID(userId)
		if err != nil {
			responseFailedMiddleware(c)
			return
		}

		c.Set("current_user", user)
	}
}

func responseFailedMiddleware(c *gin.Context) {
	response := helper.ApiResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
	c.AbortWithStatusJSON(http.StatusUnauthorized, response)
}
