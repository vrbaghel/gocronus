package server

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type ServerConfig struct {
	ServerHTTPPort            int
	ServerHTTPReadTimeout     time.Duration
	ServerHTTPWriteTimeout    time.Duration
	ServerHTTPShutdownTimeout time.Duration
}

type Server struct {
	ServerHTTPPort            int
	ServerHTTPShutdownTimeout time.Duration
	ServerInstance            *http.Server
	RoutingEngine             *gin.Engine
}

/** creates & returns a new instance of http server */
func NewServer(config ServerConfig) *Server {
	routingEngine := gin.Default()
	corsConfig := cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"POST", "GET", "DELETE", "PUT"},
		AllowHeaders:    []string{"Content-Type", "Authorization"},
		ExposeHeaders:   []string{"Access-Control-Allow-Origin"},
	}
	routingEngine.Use(cors.New(corsConfig))
	httpServer := &http.Server{
		Addr:         ":" + strconv.Itoa(config.ServerHTTPPort),
		Handler:      routingEngine,
		ReadTimeout:  config.ServerHTTPReadTimeout * time.Second,
		WriteTimeout: config.ServerHTTPWriteTimeout * time.Second,
	}

	return &Server{
		ServerHTTPPort:            config.ServerHTTPPort,
		ServerHTTPShutdownTimeout: config.ServerHTTPShutdownTimeout * time.Second,
		ServerInstance:            httpServer,
		RoutingEngine:             routingEngine,
	}
}

/** starts the http server on the given port */
func (s *Server) StartServer() {
	ctx, notifyShutdown := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	go func() {
		log.Println("http server listening on port", s.ServerHTTPPort)
		if err := s.ServerInstance.ListenAndServe(); err != nil {
			log.Fatalf("unable to register port %+v\n", err)
		}
	}()
	<-ctx.Done()
	notifyShutdown()
}

/** stops the http server gracefully */
func (s *Server) StopServer() {
	ctx, cancel := context.WithTimeout(context.Background(), s.ServerHTTPShutdownTimeout)
	defer cancel()
	if err := s.ServerInstance.Shutdown(ctx); err != nil {
		log.Fatalf("force stopped server : released port %d\t%+v", s.ServerHTTPPort, err)
	}
	log.Println("stopped http server : released port", s.ServerHTTPPort)
}
