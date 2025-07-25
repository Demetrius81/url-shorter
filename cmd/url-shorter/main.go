package main

import (
	"fmt"
	"log/slog"
	"os"
	"url-shorter/internal/config"
	"url-shorter/internal/lib/logger/sl"
	"url-shorter/internal/storage/sqlite"
	mwLogger "url-shorter/internal/http-server/middleware/logger"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {

	// TODO: init config: cleanenv VV
	cfg := config.MustLoad()

	fmt.Println(cfg)

	// TODO: init logger: slog VV

	log := setupLogger(cfg.Env)
	log.Info("starting url-shortener", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	// TODO: init storage: sqlite VV
	storage, err := sqlite.New(cfg.StoragePath)

	if err != nil {
		fmt.Println(">>>", cfg.StoragePath)
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(4)
	}

	_ = storage

	// TODO: init router: chi, "chi render" VV

	router := chi.NewRouter()

	// middleware

	router.Use(middleware.RequestID)
	// router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)


	// TODO: run server
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
