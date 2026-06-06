package middlewares

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func ValidarToken(c *gin.Context) {

	header := c.GetHeader("Authorization")

	if header == "" {
		c.JSON(400, gin.H{"error": "Token requerido"})
		c.Abort()
		return

	}
	tokenString := strings.TrimPrefix(header, "Bearer ")

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {

		return []byte(os.Getenv("JWT_SECRET")), nil

	})

	if err != nil || !token.Valid {
		c.JSON(400, gin.H{"error": "token invalido"})
		c.Abort()
		return

	}

	claims := token.Claims.(jwt.MapClaims)
	c.Set("id", uint(claims["id"].(float64)))

	c.Next()

}
