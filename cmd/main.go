package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ipxsandbox/config"
	"github.com/ipxsandbox/internal/routes"
	"github.com/ipxsandbox/internal/pkg/redis"
)

func main() {
	db := config.InitDB()
	redis.InitRedis()

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	routes.InitRoutes(r, db)

	r.Run(":8080")
}