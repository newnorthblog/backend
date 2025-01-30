package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	apiHttp "github.com/newnorthblog/backend/internal/api/http"
	"github.com/newnorthblog/backend/internal/config"
	"github.com/newnorthblog/backend/internal/db"
	"github.com/newnorthblog/backend/internal/pkg/tokenmanager"
	"github.com/newnorthblog/backend/internal/repository"
	"github.com/newnorthblog/backend/internal/server"
	"github.com/newnorthblog/backend/internal/service"
	"github.com/newnorthblog/backend/pkg/logger"
)

func main() {
	// Init config
	cfg := config.MustLoad()

	// Init logger
	logger := logger.SetupLogger(cfg.Env)
	logger.Info("start new north blog backend",
		"env", cfg.Env,
	)
	logger.Debug("debug messages are enabled")

	// Init database
	dbPostgres, err := db.New(cfg.Database)
	if err != nil {
		logger.Error("postgres connection problem", "error", err)
		os.Exit(1)
	}
	defer func() {
		err = dbPostgres.Close()
		if err != nil {
			logger.Error("error when closing", "error", err)
		}
	}()
	logger.Info("postgres connection done")

	// Init services, repositories, handlers
	repos := repository.NewRepositories(dbPostgres)
	tokenManager, err := tokenmanager.NewManager(cfg.JWT.SecretKey, cfg.JWT.AccessTokenTTL, cfg.JWT.RefreshTokenTTL)
	if err != nil {
		logger.Error("token manager error", "error", err)
		os.Exit(1)
	}
	services := service.NewServices(service.Deps{
		Logger:       logger,
		Config:       cfg,
		Repos:        repos,
		TokenManager: tokenManager,
	})
	handlers := apiHttp.NewHandlers(
		services,
		logger,
		tokenManager,
	)

	// Init HTTP server
	srv := server.NewServer(cfg, handlers.Init(cfg))
	go func() {
		if err := srv.Run(); !errors.Is(err, http.ErrServerClosed) {
			logger.Error("error occurred while running http server", "error", err)
		}
	}()
	logger.Info("server started")

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	const timeout = 5 * time.Second

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err := srv.Stop(ctx); err != nil {
		logger.Error("failed to stop server", "error", err)
	}

	logger.Info("app stopped")
}
