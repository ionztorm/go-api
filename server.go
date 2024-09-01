package main

import (
	"github.com/LeonLonsdale/go-web-api/db"
	"github.com/LeonLonsdale/go-web-api/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDB()
	server := gin.Default()

	routes.RegisterRoutes(server)
	server.Run(":8080")
}
