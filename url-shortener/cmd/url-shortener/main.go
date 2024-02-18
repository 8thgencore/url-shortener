package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	ssogrpc "url-shortener/internal/clients/sso/grpc"
	"url-shortener/internal/config"
	"url-shortener/internal/http-server/handlers/greeting"
	"url-shortener/internal/http-server/handlers/redirect"
	hDelete "url-shortener/internal/http-server/handlers/url/delete"
	"url-shortener/internal/http-server/handlers/url/save"
	mwLogger "url-shortener/internal/http-server/middleware/logger"
	"url-shortener/internal/lib/logger/handlers/slogpretty"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/storage/sqlite"

	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// Load configuration
	cfg := config.MustLoad()

	// Set up logger based on configuration
	log := newSlogLogger(cfg.Log.Slog)

	// Log information about the start of the application
	log.Info("starting url-shortener", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	// Initial gRPC sso client
	ssoClient, err := ssogrpc.New(
		context.Background(),
		log,
		cfg.Clients.SSO.Address,
		cfg.Clients.SSO.Timeout,
		cfg.Clients.SSO.RetriesCount,
	)
	if err != nil {
		log.Error("failed to init sso client", sl.Err(err))
		os.Exit(1)
	}

	// Example used
	isAdmin, _ := ssoClient.IsAdmin(context.Background(), 1)
	log.Info("is_admin", slog.Bool("is_admin", isAdmin))

	// Initialize storage (SQLite in this case)
	dirPath := filepath.Dir(cfg.StoragePath)
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
			log.Error("failed to make storage dirs", sl.Err(err))
			os.Exit(1)
		}
	}
	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("Failed to init storage", sl.Err(err))
		os.Exit(1)
	}
	defer func() {
		if err := storage.Close(); err != nil {
			log.Error("failed to close storage", sl.Err(err))
		}
	}()

	// Create a new Chi router
	router := chi.NewRouter()

	// Middleware setup
	router.Use(middleware.RequestID)
	// TODO: Log requests with mismatched route using slog instead of chi logger
	router.Use(middleware.Logger)
	router.Use(mwLogger.New(log))
	// TODO: use slog to log panic and github.com/maruel/panicparse to handle it
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	// Create a new FileServer to serve static files from the "./static" directory
	fs := http.FileServer(http.Dir("./static"))

	// Handle requests to URLs starting with "/static/" by stripping the prefix and serving files from the file server
	router.Handle("/static/*", http.StripPrefix("/static/", fs))

	// Define a route for "/url" with basic authentication
	router.Route("/url", func(r chi.Router) {
		r.Use(middleware.BasicAuth("url-shortener", map[string]string{
			cfg.HttpServer.User: cfg.HttpServer.Password,
		}))
		// Uncomment and customize the following lines based on your routes
		// r.Post("/", save.New(log, storage))
		// r.Delete("/{alias}", hDelete.New(log, storage))
	})

	// Define routes for saving, deleting, and redirecting
	router.Get("/", greeting.New(log, "./static"))
	router.Post("/", save.New(log, storage))
	router.Delete("/{alias}", hDelete.New(log, storage))
	router.Get("/{alias}", redirect.New(log, storage))

	// Log information about the server start
	log.Info("starting server", slog.String("address", cfg.Address))

	// Set up signals for graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Create an HTTP server with specified configurations
	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HttpServer.Timeout,
		WriteTimeout: cfg.HttpServer.Timeout,
		IdleTimeout:  cfg.HttpServer.Timeout,
	}

	// Start the server in a goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error("failed to start Server", slog.String("error", err.Error()))
		}
	}()

	// Log information about the server start
	log.Info("server started")

	// Wait for a signal to stop the server
	<-done
	log.Info("stopping server")

	// Set up a context with a timeout for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), cfg.HttpServer.ShutdownTimeout)
	defer cancel()

	// Attempt to gracefully stop the server
	if err := srv.Shutdown(ctx); err != nil {
		log.Error("failed to stop server", sl.Err(err))
		return
	}

	// Log information about the server stop
	log.Error("server stopped")
}

func newSlogLogger(c config.Slog) *slog.Logger {
	o := &slog.HandlerOptions{Level: c.Level, AddSource: c.AddSource}
	w := os.Stdout
	var h slog.Handler

	switch c.Format {
	case "pretty":
		h = slogpretty.NewHandler().
			WithAddSource(c.AddSource).
			WithLevel(c.Level).
			WithLevelEmoji(c.Pretty.Emoji).
			WithFieldsFormat(c.Pretty.FieldsFormat)
	case "json":
		h = slog.NewJSONHandler(w, o)
	case "text":
		h = slog.NewTextHandler(w, o)
	}
	return slog.New(h)
}
