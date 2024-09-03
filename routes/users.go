package routes

import (
	"net/http"

	"github.com/LeonLonsdale/go-web-api/models"
	"github.com/LeonLonsdale/go-web-api/utils"
	"github.com/gin-gonic/gin"
)

func signup(context *gin.Context) {
	var user models.User

	err := context.ShouldBindJSON(&user)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "could not parse the submitted data"})
		return
	}

	err = user.Save()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "a problem occurred while saving the user", "error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "signup successful"})
}

func login(context *gin.Context) {
	var user models.User

	err := context.ShouldBindJSON(&user)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "could not parse the submitted details"})
		return
	}

	err = user.ValidateCredentials()
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	token, err := utils.GenrateToken(user.Email, user.ID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not authenticate user", "error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "logged in", "token": token})
}
