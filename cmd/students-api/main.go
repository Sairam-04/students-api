package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Sairam-04/students-api/pkg/config"
)

func main() {
	// load config
	cfg := config.MustLoad()
	// database setup

	// setup router
	router := http.NewServeMux()
	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})
	// setup server
	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}
	slog.Info("server running at", slog.String("address", cfg.Addr) )
	done := make(chan os.Signal, 1)
	// signal.notify method sends signal in done channel and code is unblocked and <- done below code runs
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	// if any interup signal make it non blocking ListenAndServer is blocking

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("failed to start server")
		}
	}()
	<-done

	// logic for server stopping
	slog.Info("shutting down the server")
	// server.shutdown takes time so we need to block the requests coming
	// use context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := server.Shutdown(ctx)
	if err != nil {
		slog.Error("failed to shutdown server", slog.String("error", err.Error()))
	}
	slog.Info("server shutdown successfully")
}
