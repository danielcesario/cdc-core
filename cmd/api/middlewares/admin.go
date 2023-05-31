package middlewares

import (
	"github.com/gin-gonic/gin"
)

func IsSuperAdmin() gin.HandlerFunc {
	return func(context *gin.Context) {
		roles := context.MustGet("roles")

		var isSuperAdmin bool = false
		for _, role := range roles.([]string) {
			if role == "SUPER_ADMIN" {
				isSuperAdmin = true
			}
		}

		if !isSuperAdmin {
			context.JSON(403, gin.H{"error": "you have no roles to manage this application"})
			context.Abort()
		}

		context.Next()
	}
}
