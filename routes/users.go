package routes

import (
	"net/http"

	"github.com/LeonLonsdale/go-web-api/models"
	"github.com/gin-gonic/gin"
)

func signup(context *gin.Context) {
	var user models.User

	error := context.ShouldBindJSON(&user)

	if error != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "could not parse the submitted data"})
		return
	}

	error = user.Save()

	if error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "a problem occurred while saving the user"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "signup successful"})
}

func login(context *gin.Context) {
	var user models.User

	error := context.ShouldBindJSON(&user)

	if error != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "could not parse the submitted details"})
		return
	}

	error = user.ValidateCredentials()
	if error != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"message": error.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "logged in"})
}
