package api

import (
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	listenAddr string
	router     *gin.Engine
	hub        *Hub
}

var serv *Server

func GetServer() *Server {
	return serv
}

func NewServer(listenAddr string) *Server {
	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowMethods = []string{"GET", "POST"}
	config.AllowHeaders = []string{"Authorization", "Content-Type", "Origin"}

	r.Use(cors.New(config))

	serv = &Server{
		listenAddr: listenAddr,
		router:     r,
		hub:        NewHub(),
	}

	return serv
}

func (s *Server) Start() error {

	s.router.GET("ws", s.ChatHandler)

	s.router.Run(s.listenAddr)
	return http.ListenAndServe(s.listenAddr, nil)
}

func (s *Server) ChatHandler(c *gin.Context) {
	log.Println("new connection !")
	err := Handler(c.Writer, c.Copy().Request, s.hub)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(500)
	}
	log.Println("connection closed !")
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
