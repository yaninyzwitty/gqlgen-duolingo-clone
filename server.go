package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/yaninyzwitty/gqlgen-duolingo-clone/graph"
	"github.com/yaninyzwitty/gqlgen-duolingo-clone/internal/database"
	"github.com/yaninyzwitty/gqlgen-duolingo-clone/pkg"
)

var (
	cfg      pkg.Config
	password string
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	file, err := os.Open("config.yaml")
	if err != nil {
		slog.Error("failed to open config.yaml", "error", err)
		os.Exit(1)
	}
	defer file.Close()

	if err = cfg.LoadConfig(file); err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	if s := os.Getenv("DB_PASSWORD"); s != "" {
		password = s
	}

	dbConfig := database.DbConfig{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: password,
		DbName:   cfg.Database.Database,
		MaxConn:  cfg.Database.MaxRetries,
	}

	pool, err := dbConfig.NewPgxPool(ctx, cfg.Database.MaxRetries)
	if err != nil {
		slog.Error("failed to create database pool", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	err = dbConfig.Ping(ctx, pool, cfg.Database.MaxRetries)
	if err != nil {
		slog.Error("failed to ping database", "error", err)
		os.Exit(1)
	}

	mux := chi.NewRouter()
	mux.Use(middleware.Logger)

	resolvers := &graph.Resolver{
		Pool: pool,
	}
	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: resolvers}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	mux.Handle("/", playground.Handler("GraphQL playground", "/query"))
	mux.Handle("/query", srv)
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: mux,
	}

	stopCH := make(chan os.Signal, 1)
	signal.Notify(stopCH, os.Interrupt, syscall.SIGTERM)

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	go func() {
		slog.Info("SERVER starting", "port", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("failed to start server", "error", err)
			os.Exit(1)
		}
	}()

	<-stopCH
	slog.Info("shutting down the server...")
	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("failed to shutdown server", "error", err)
		os.Exit(1)
	} else {
		slog.Info("server stopped gracefully")
	}

}
