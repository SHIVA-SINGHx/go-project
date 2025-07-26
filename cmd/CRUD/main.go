package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/SHIVA-SINGHx/Go-Project/internal/config"
	"github.com/SHIVA-SINGHx/Go-Project/internal/http/handlers"

)

func main() {
	cfg := config.MustLoad()

	mux := http.NewServeMux()

	// Home route
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Write([]byte("🚀 Welcome To My Go CRUD Project"))
	})

	// Routes
	mux.HandleFunc("/api/student", handlers.CreateStudent())               // POST METHOOD
	mux.HandleFunc("/api/student/", handlers.StudentHandlerWithID()) // GET, PUT, DELETE METHOOD

	// Server setup
	server := http.Server{
		Addr:    cfg.HTTPServer.Addr,
		Handler: mux,
	}

	slog.Info("server started", slog.String("address", cfg.HTTPServer.Addr))
	fmt.Printf("Server started on %s\n", cfg.HTTPServer.Addr)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("Server error: %s", err)
		}
	}()

	<-done
	slog.Info("Shutting Down The Server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Error shutting down", slog.String("error", err.Error()))
	}
}
