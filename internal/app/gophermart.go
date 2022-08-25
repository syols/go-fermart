package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/syols/go-devops/config"
	"github.com/syols/go-devops/internal/handlers"
	"github.com/syols/go-devops/internal/pkg"
)

type Server struct {
	server   http.Server
	settings config.Config
}

func NewServer(settings config.Config) (Server, error) {
	auth := pkg.NewAuthorizer(settings)
	db, err := pkg.NewDatabaseConnection(settings)
	if err != nil {
		return Server{}, err
	}

	sess, err := pkg.NewSession()
	if err != nil {
		log.Print(err.Error())
	}

	return Server{
		server: http.Server{
			Addr:    settings.ServerAddress.String(),
			Handler: router(db, auth, sess),
		},
		settings: settings,
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

func router(db pkg.Database, auth pkg.Authorizer, sess *pkg.Session) *gin.Engine {
	router := gin.Default()
	router.Use(gin.Recovery())
	router.Use(handlers.LoggerMiddleware())
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	router.GET("/healthcheck", handlers.Healthcheck)

	api := router.Group("/api/user")
	api.POST("/register", handlers.Register(db, auth))
	api.POST("/login", handlers.Login(db, auth))

	authorized := api.Group("/")
	authorized.Use(handlers.AuthMiddleware(db, auth))

	orders := authorized.Group("/")
	orders.POST("/orders", handlers.CreatePurchase(db, sess))
	orders.GET("/orders", handlers.Purchases(db))

	balance := authorized.Group("/")
	balance.GET("/balance", handlers.Balance(db))
	balance.POST("/balance/withdraw", handlers.CreateWithdraw(db))
	balance.GET("/withdrawals", handlers.Withdrawals(db))
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
