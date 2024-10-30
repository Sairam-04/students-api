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

	"github.com/Sairam-04/students-api/internal/config"
	"github.com/Sairam-04/students-api/internal/http/handlers/student"
	"github.com/Sairam-04/students-api/internal/storage/postgres"
	_ "github.com/Sairam-04/students-api/internal/storage/sqlite"
)

func main() {
	// load config
	cfg := config.MustLoad()

	// database setup
	storage, err := postgres.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	slog.Info("storage initialized", slog.String("env", cfg.Env), slog.String("version", "1.0.0"))

	// setup router
	router := http.NewServeMux()
	router.HandleFunc("POST /api/students", student.New(storage))
	router.HandleFunc("GET /api/students/{id}", student.GetByID(storage))
	router.HandleFunc("GET /api/students", student.GetList(storage))
	router.HandleFunc("PUT /api/students/{id}", student.UpdateStudent(storage))
	router.HandleFunc("DELETE /api/students/{id}", student.DeleteStudent(storage))
	// setup server
	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}
	slog.Info("server running at", slog.String("address", cfg.Addr))
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
	err = server.Shutdown(ctx)
	if err != nil {
		slog.Error("failed to shutdown server", slog.String("error", err.Error()))
	}
	slog.Info("server shutdown successfully")
}
