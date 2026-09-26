package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gtantech/p2/internal/models"
	"github.com/gtantech/p2/internal/server"
)

func gracefulShutdown(server *server.Server, done chan bool) {
	apiServer := server.Server
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Listen for the interrupt signal.
	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")
	stop() // Allow Ctrl+C to force shutdown

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := apiServer.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	log.Println("Server exiting")
	log.Println("Closing database...")
	if err := server.CloseDb(); err != nil {
		log.Printf("Database closed with error: %v", err)
	} else {
		log.Println("Database closed successfully")
	}

	// Notify the main goroutine that the shutdown is complete
	done <- true
}

func main() {

	port, _ := strconv.Atoi(os.Getenv("PORT"))
	serverConfig := models.ServerConfig{Port: port}
	server := server.NewServer(serverConfig)

	// Create a done channel to signal when the shutdown is complete
	done := make(chan bool, 1)

	// Run graceful shutdown in a separate goroutine
	go gracefulShutdown(server, done)
	log.Printf("Listening on http://localhost:%d", serverConfig.Port)
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("http server error: %s", err))
	}
	// Wait for the graceful shutdown to complete
	<-done
	log.Println("Graceful shutdown complete.")
}
