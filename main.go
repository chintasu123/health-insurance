package main

import (
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	server := gin.New()
	err := server.Run()
	if err != nil {
		log.Println("unable to server :", err)
	}
}
