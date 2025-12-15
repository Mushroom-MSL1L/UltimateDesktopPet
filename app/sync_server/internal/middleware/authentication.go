package middleware

import (
	"github.com/Mushroom-MSL1L/UltimateDesktopPet/app/sync_server/internal/controllers"
	"github.com/Mushroom-MSL1L/UltimateDesktopPet/app/sync_server/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		token, err := jwt.Parse(tokenString,
			func(token *jwt.Token) (interface{}, error) {
				return []byte("secret"), nil
			})
		if err != nil {
			c.JSON(401, controllers.ResponseMessage{Message: "token error"})
			c.Abort()
			return
		}
		claims := token.Claims.(jwt.MapClaims)
		username := claims["username"].(string)

		if !models.ValidateUsername(username) {
			c.JSON(401, controllers.ResponseMessage{Message: "Unauthorized"})
			c.Abort()
			return
		}
		c.Set("username", username)
		c.Next()
	}
}
