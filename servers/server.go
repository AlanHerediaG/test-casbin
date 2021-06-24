package servers

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/AlanHerediaG/test-casbin/auth"

	fileadapter "github.com/casbin/casbin/persist/file-adapter"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/joho/godotenv"
)

type Server struct {
	Router      *gin.Engine
	FileAdapter *fileadapter.Adapter
	RedisCli    *redis.Client
	RD          auth.IAuth
	TK          auth.IToken
}

var HttpServer Server

func (server *Server) Initialize(redis_host, redis_port, redis_password string) {
	server.Router = gin.Default()
	server.RedisCli = NewRedisDB(redis_host, redis_port, redis_password)
	server.FileAdapter = fileadapter.NewAdapter("config/basic_policy.csv")
	server.RD = auth.NewAuthService(server.RedisCli)
	server.InitializeRoutes()
}

func NewRedisDB(host, port, password string) *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: password,
		DB:       0,
	})

	return redisClient
}

func (server *Server) Run(addr string) {
	fmt.Printf("Listen on port %s \n", addr)
	log.Fatal(http.ListenAndServe(addr, server.Router))
}

func Run() {
	HttpServer = Server{}
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting env, not comming through %v", err)
	} else {
		fmt.Println("We are getting the env values")
	}

	appAddr := ":" + os.Getenv("PORT")
	redis_host := os.Getenv("REDIS_HOST")
	redis_port := os.Getenv("REDIS_PORT")
	redis_password := os.Getenv("REDIS_PASSWORD")

	HttpServer.Initialize(redis_host, redis_port, redis_password)

	HttpServer.Run(appAddr)
}

func (server *Server) Close() {
	if server.RedisCli != nil {
		if err := server.RedisCli.Close(); err != nil {
			log.Fatal(err)
		}
		server.RedisCli = nil
	}
}
