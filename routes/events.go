package routes

import (
	"net/http"
	"strconv"

	"github.com/LeonLonsdale/go-web-api/models"
	"github.com/gin-gonic/gin"
)

func getEvent(context *gin.Context) {
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "an event ID is required"})
		return
	}

	event, err := models.GetEventByID(eventId)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not retrieve event"})
		return
	}

	context.JSON(http.StatusOK, event)
}

func getEvents(context *gin.Context) {
	events, err := models.GetAllEvents()

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not retrieve events"})
		return
	}

	context.JSON(http.StatusOK, events)

}

func createEvent(context *gin.Context) {
	var event models.Event

	err := context.ShouldBindJSON(&event) // store the data from request body in event var
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "could not parse submitted data", "error": err})
		return
	}

	event.ID = 1
	event.UserID = 1

	err = event.Save()

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "unable to save event"})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "event created", "event": event})
}
