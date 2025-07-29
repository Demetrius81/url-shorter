package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"url-shorter/internal/config"
	"url-shorter/internal/http-server/handlers/url/redirect"
	"url-shorter/internal/http-server/handlers/url/remove"
	"url-shorter/internal/http-server/handlers/url/save"
	mwLogger "url-shorter/internal/http-server/middleware/logger"
	"url-shorter/internal/lib/logger/handlers/slogpretty"
	"url-shorter/internal/lib/logger/sl"
	"url-shorter/internal/storage/sqlite"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {

	// init config: cleanenv
	cfg := config.MustLoad()

	// init logger: slog

	log := setupLogger(cfg.Env)
	log.Info("starting url-shortener", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	// init storage: sqlite
	storage, err := sqlite.New(cfg.StoragePath)

	if err != nil {
		fmt.Println(">>>", cfg.StoragePath)
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(4)
	}

	// _ = storage

	// init router: chi, "chi render"

	router := chi.NewRouter()

	// middleware

	router.Use(middleware.RequestID)
	// router.Use(middleware.RealIP)
	// router.Use(middleware.Logger)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Route("/url", func(r chi.Router) {
		r.Use(middleware.BasicAuth("url-shorter", map[string]string{
			cfg.HTTPServer.User: cfg.HTTPServer.Password,
		}))

		r.Post("/url", save.New(log, storage))
		r.Delete("/url/{alias}", remove.New(log, storage))

	})

	router.Get("/{alias}", redirect.New(log, storage))

	log.Info("starting server", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Error("server stopped")

	// TODO: run server
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettyLogger()
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}

func setupPrettyLogger() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
