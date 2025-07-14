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

	"github.com/codingsher/golang-student-api/internal/config"
	"github.com/codingsher/golang-student-api/internal/http/handlers/student"
	"github.com/codingsher/golang-student-api/internal/storage/sqlite"
)

func main() {
	cfg := config.MustLoad()

	storage, err := sqlite.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	slog.Info("storage initializated", slog.String("env", cfg.Env), slog.String("version", "1.0.0"))

	router := http.NewServeMux()

	router.HandleFunc("POST /api/students", student.New(storage))
	router.HandleFunc("GET /api/students/{id}", student.GetByID(storage))
	router.HandleFunc("GET /api/students", student.GetList(storage))
	router.HandleFunc("PATCH /api/students/{id}", student.UpdateByID(storage))
	router.HandleFunc("DELETE /api/students/{id}", student.DeleteByID(storage))

	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}
	slog.Info("server started", slog.String("address", cfg.Addr))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("server failed to start!")
		}
	}()

	<-done

	slog.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		slog.Error("serve failed to shutdown", slog.String("error", err.Error()))
	}

	slog.Info("server shutdown success...")
}
