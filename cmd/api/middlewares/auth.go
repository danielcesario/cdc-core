package middlewares

import (
	"github.com/danielcesario/cdc-core/internal/auth"
	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(context *gin.Context) {
		tokenString := context.GetHeader("Authorization")
		if tokenString == "" {
			context.JSON(401, gin.H{"error": "request does not contain an access token"})
			context.Abort()
			return
		}

		claim, err := auth.ValidateToken(tokenString)
		if err != nil {
			context.JSON(401, gin.H{"error": err.Error()})
			context.Abort()
			return
		}

		context.Set("user_code", claim.Code)
		context.Set("user_email", claim.Email)
		context.Set("roles", claim.Roles)
		context.Next()
	}
}
