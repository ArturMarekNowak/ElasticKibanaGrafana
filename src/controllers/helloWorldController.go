package controllers

import (
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"time"
)

func GetHelloWorld(c *gin.Context) {
	sleep(1000)
	c.JSON(http.StatusOK, gin.H{"message": "Hello World!"})
}

func sleep(ms int) {
	now := time.Now()
	n := rand.Intn(ms + now.Second())
	time.Sleep(time.Duration(n) * time.Millisecond)
}
