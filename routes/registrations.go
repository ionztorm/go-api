package routes

import (
	"net/http"
	"strconv"

	"github.com/LeonLonsdale/go-web-api/models"
	"github.com/gin-gonic/gin"
)

func makeRegistration(context *gin.Context) {
	userId := context.GetInt64("userId")
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "could not parse event ID"})
		return
	}

	event, err := models.GetEventByID(eventId)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "no event found for that ID"})
		return
	}

	isAlreadyRegistered, err := models.CheckIfAlreadyRegistered(userId, eventId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "an error occurred checking for existing registrations", "error": err.Error()})
		return
	}

	if isAlreadyRegistered {
		context.JSON(http.StatusConflict, gin.H{"message": "You have already registered for this event"})
		return
	}

	err = event.Register(userId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "unable to register for event", "error": err.Error()})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "registered"})
}

func cancelRegistration(context *gin.Context) {
	userId := context.GetInt64("userId")
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "could not parse event ID"})
		return
	}

	var event models.Event
	event.ID = eventId
	err = event.CancelRegistration(userId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "unable to delete this registration at this time", "error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
