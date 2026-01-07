// cmd/server/main.go
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/SecureParadise/go_attendence/internal/api/routes"
	"github.com/SecureParadise/go_attendence/internal/config"
	"github.com/SecureParadise/go_attendence/internal/db"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Signals that will trigger graceful shutdown
var interruptSignals = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGINT,
}

func main() {
	// --------------------------------------------------
	// 1️⃣ Load configuration using Viper
	// --------------------------------------------------
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	// --------------------------------------------------
	// 2️⃣ Connect to PostgreSQL using pgxpool
	// --------------------------------------------------
	ctx := context.Background()
	connPool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal("cannot connect to database:", err)
	}
	defer connPool.Close()

	// Test database connection
	if err := connPool.Ping(ctx); err != nil {
		log.Fatal("cannot ping database:", err)
	}
	log.Println("connected to database")

	// --------------------------------------------------
	// 3️⃣ Create store and server
	// --------------------------------------------------
	store := db.NewStore(connPool)
	server, err := routes.NewServer(cfg, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	// --------------------------------------------------
	// 4️⃣ Create HTTP server
	// --------------------------------------------------
	httpServer := &http.Server{
		Addr:    cfg.HTTPServerAddress,
		Handler: server.GetRouter(),
	}

	// --------------------------------------------------
	// 5️⃣ Channel to listen for OS signals
	// --------------------------------------------------
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, interruptSignals...)

	// --------------------------------------------------
	// 6️⃣ Start server in a goroutine
	// --------------------------------------------------
	go func() {
		log.Println("HTTP server started on", cfg.HTTPServerAddress)

		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("server error:", err)
		}
	}()

	// --------------------------------------------------
	// 7️⃣ Wait for shutdown signal
	// --------------------------------------------------
	<-quit
	log.Println("shutdown signal received")

	// --------------------------------------------------
	// 8️⃣ Create context with timeout for graceful shutdown
	// --------------------------------------------------
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// --------------------------------------------------
	// 9️⃣ Shutdown HTTP server gracefully
	// --------------------------------------------------
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Println("server forced to shutdown:", err)
	} else {
		log.Println("server shutdown completed gracefully")
	}

	log.Println("application exited cleanly")
}
