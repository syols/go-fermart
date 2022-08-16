package app

import (
	"context"
	"github.com/gin-gonic/contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/syols/go-devops/config"
	"github.com/syols/go-devops/internal/database"
	"github.com/syols/go-devops/internal/handlers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type Server struct {
	server   http.Server
	database database.Database
	settings config.Config
}

func NewServer(settings config.Config) (Server, error) {
	db, err := database.NewDatabase(settings.DatabaseConnectionString)
	if err != nil {
		return Server{}, err
	}

	return Server{
		server: http.Server{
			Addr:    settings.Address(),
			Handler: router(),
		},
		settings: settings,
		database: db,
	}, nil
}

func (s *Server) Run() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	s.shutdown(ctx)

	if err := s.server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func router() *gin.Engine {
	router := gin.Default()
	router.Use(gin.Recovery())
	router.Use(gzip.Gzip(gzip.DefaultCompression))

	router.GET("/healthcheck", handlers.Healthcheck)
	api := router.Group("/api/user")
	api.POST("/register", handlers.Register)
	api.POST("/login", handlers.Login)
	api.POST("/orders", handlers.SetUserOrders)
	api.GET("/orders", handlers.UserOrders)

	balance := api.Group("/balance")
	balance.GET("/", handlers.Balance)
	balance.GET("/withdraw", handlers.Withdraw)
	balance.GET("/withdrawals", handlers.Withdrawals)

	return router
}

func (s *Server) shutdown(ctx context.Context) {
	go func() {
		<-ctx.Done()
		if err := s.server.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
}
