package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"

	"github.com/thangchung/go-coffeeshop/cmd/auth/config"
	handler "github.com/thangchung/go-coffeeshop/internal/auth/handlers"
	repository "github.com/thangchung/go-coffeeshop/internal/auth/repo"
	"github.com/thangchung/go-coffeeshop/internal/auth/service"
)

func main() {
	// ── Config ──────────────────────────────────────────────────────────────
	cfg, err := config.NewConfig()
	if err != nil {
		logrus.WithError(err).Fatal("failed to load config")
	}

	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)

	logrus.WithFields(logrus.Fields{
		"name":    cfg.Name,
		"version": cfg.Version,
	}).Info("⚡ starting auth-service")

	// ── Database ─────────────────────────────────────────────────────────────
	db, err := sql.Open("postgres", cfg.Postgres.URL)
	if err != nil {
		logrus.WithError(err).Fatal("failed to open db connection")
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		logrus.WithError(err).Fatal("failed to ping database")
	}
	logrus.Info("✅ database connected")

	// ── Repository / Service / Handler ───────────────────────────────────────
	userRepo := repository.NewUserRepository(db)
	if err = userRepo.Migrate(); err != nil {
		logrus.WithError(err).Fatal("failed to run migrations")
	}

	authSvc := service.NewAuthService(userRepo, cfg.JWT.Secret, cfg.JWT.ExpireHours)
	authHandler := handler.NewAuthHandler(authSvc)

	// ── Router ───────────────────────────────────────────────────────────────
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(requestLogger())

	api := router.Group("/auth")
	authHandler.RegisterRoutes(api)

	// ── Server ───────────────────────────────────────────────────────────────
	addr := fmt.Sprintf("%s:%d", cfg.HTTP.Host, cfg.HTTP.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	go func() {
		logrus.WithField("address", addr).Info("🌏 auth-service listening")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.WithError(err).Fatal("server error")
		}
	}()

	// ── Graceful shutdown ────────────────────────────────────────────────────
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	logrus.Info("shutting down auth-service…")
}

// requestLogger is a simple Gin middleware that logs every request with logrus.
func requestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		logrus.WithFields(logrus.Fields{
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
			"status": c.Writer.Status(),
		}).Info("request")
	}
}
