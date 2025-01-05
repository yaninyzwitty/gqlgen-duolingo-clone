package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/yaninyzwitty/gqlgen-duolingo-clone/graph"
	"github.com/yaninyzwitty/gqlgen-duolingo-clone/internal/pkg"
)

var cfg pkg.Config

func main() {
	loadConfig()

	router := setupRouter()

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: router,
	}

	startServer(server)
}

func loadConfig() {
	file, err := os.Open("config.yaml")
	if err != nil {
		slog.Error("failed to open config.yaml", "error", err)
		os.Exit(1)
	}
	defer file.Close()

	if err := cfg.LoadConfig(file); err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}
}

func setupRouter() http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))
	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	return router
}

func startServer(server *http.Server) {
	stopCH := make(chan os.Signal, 1)
	signal.Notify(stopCH, os.Interrupt, syscall.SIGTERM)

	go func() {
		slog.Info("server is listening on :" + fmt.Sprintf("%d", cfg.Server.Port))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("failed to start server", "error", err)
			os.Exit(1)
		}
	}()

	<-stopCH

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.Timeout)
	defer cancel()

	slog.Info("shutting down the server...")
	if err := server.Shutdown(ctx); err != nil {
		slog.Error("failed to shutdown server", "error", err)
		os.Exit(1)
	} else {
		slog.Info("server stopped gracefully")
	}
}
