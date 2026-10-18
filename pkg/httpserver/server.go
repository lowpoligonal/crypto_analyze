package httpserver

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ChipherController interface {
	Init(r *gin.RouterGroup)
}

type Server struct {
	httpServer *http.Server
}

func NewServer(controller ChipherController) *Server {
	router := gin.Default()

	router.Static("/static", "./web/static")
	router.GET("/", func(c *gin.Context) {
		c.File("./web/static/index.html")
	})

	api := router.Group("/api")
	controller.Init(api)

	return &Server{
		httpServer: &http.Server{
			Addr:    ":8080",
			Handler: router,
		},
	}
}

func (s *Server) Run(addr ...string) error {
	if len(addr) > 0 {
		s.httpServer.Addr = addr[0]
	}
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
