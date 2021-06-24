// References
// https://tienbm90.medium.com/authentication-and-authorization-in-gin-application-with-jwt-and-casbin-a56bbbdec90b

package main

import (
	"log"

	"github.com/AlanHerediaG/test-casbin/servers"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var (
	router *gin.Engine
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	servers.Run()
	log.Println("Server exiting")
}
