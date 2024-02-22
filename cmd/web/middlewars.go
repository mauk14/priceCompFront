package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

type envelope map[string]any

func ErrorResponse(c *gin.Context, status int, message any) {
	env := envelope{"error": message}
	fmt.Println(message)
	c.IndentedJSON(status, env)
}

func recoverPanic() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.Header("Connection", "close")
				c.IndentedJSON(http.StatusInternalServerError, gin.H{
					"error": fmt.Errorf("%s", err).Error(),
				})
			}
		}()
		c.Next()
	}
}

func requireAuth(c *gin.Context) {
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		c.Next()
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header)
		}
		return []byte("chertMama"), nil
	})

	if err != nil {
		ErrorResponse(c, http.StatusUnauthorized, "invalid authentication credentials")
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			ErrorResponse(c, http.StatusUnauthorized, "invalid authentication credentials")
			return
		}

		//user, err := app.userManager.UserInfo(context.Background(), claims["sub"].(string))
		//if err != nil {
		//	ErrorResponse(c, http.StatusUnauthorized, "invalid authentication credentials")
		//	return
		//}

		name := claims["name"].(string)
		email := claims["email"].(string)
		userId := int64(claims["userId"].(float64))
		c.Set("name", name)
		c.Set("email", email)
		c.Set("userId", userId)

		c.Next()
	}
}

type userToken struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
