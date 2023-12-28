package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/brianeno/gokafka/db"
)

var (
	ListenAddr = "localhost:8080"
	RedisAddr  = "localhost:6379" // modify as needed
)

func main() {
	redis, err := db.NewDatabase(RedisAddr)
	if err != nil {
		log.Fatalf("Failed to connect to redis: %s", err.Error())
	}

	router := initRouter(redis)
	router.Run(ListenAddr)
}

func initRouter(database *db.RedisDatabase) *gin.Engine {
	r := gin.Default()
	r.GET("/storypoints/:title", func(c *gin.Context) {
		title := c.Param("title")
		project, err := database.GetProject(title)
		if err != nil {
			if err == db.ErrNil {
				c.JSON(http.StatusNotFound, gin.H{"error": "No project record found for " + title})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"project": project})
	})

	r.POST("/storypoints", func(c *gin.Context) {
		var userJson db.Project
		if err := c.ShouldBindJSON(&userJson); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		err := database.SaveProject(&userJson)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"user": userJson})
	})

	r.GET("/dashboard", func(c *gin.Context) {
		dashboard, err := database.GetDashboard()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"dashboard": dashboard})
	})

	return r
}
