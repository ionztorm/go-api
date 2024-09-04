package middlewares

import (
	"net/http"

	"github.com/LeonLonsdale/go-web-api/utils"
	"github.com/gin-gonic/gin"
)

func Authenticate(context *gin.Context) {
	token := context.Request.Header.Get("Authorization")
	if token == "" {
		// aborts the current response generation, and sends a response
		// without this, the other request handlers would still run.
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "You must be logged in"})
		return
	}

	UserID, err := utils.VerifyToken(token)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "You must be logged in"})
		return
	}

	context.Set("UserID", UserID)

	// proceed to the next handlers
	context.Next()
}
