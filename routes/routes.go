package routes

import (
	"github.com/LeonLonsdale/go-web-api/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	server.GET("/events", getEvents)

	// Protected Group
	protected := server.Group("/")
	protected.Use(middlewares.Authenticate)
	protected.POST("/events", createEvent)
	protected.GET("/events/:id", getEvent)
	protected.PUT("/events/:id", updateEvent)
	protected.DELETE("/events/:id", deleteEvent)
	protected.POST("/events/:id/registers", makeRegistration)
	protected.DELETE("/events/:id/registers", cancelRegistration)

	server.POST("/signup", signup)
	server.POST("/login", login)
}
