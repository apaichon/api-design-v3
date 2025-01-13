package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/spf13/viper"

	"api/config"
	"api/internal/auth"
	"api/internal/middleware"
	"api/internal/monitoring"
	"api/internal/payment"
)

var cfg *config.Config

func init() {
	// Load configuration
	cfg = config.NewConfig()
}

func main() {
	shutdown, err := monitoring.InitTracer(viper.GetString("TRACE_EXPORTER_URL"))
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Fatal("failed to shutdown TracerProvider: %w", err)
		}
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("/api/role", auth.CreateRoleHandler)
	mux.HandleFunc("/api/login", auth.LoginHandler)
	mux.HandleFunc("/api/register", auth.RegisterHandler)
	mux.HandleFunc("/api/logout", auth.LogoutHandler)
	mux.HandleFunc("/api/users", auth.GetUsersHandler)
	mux.HandleFunc("/api/payments", payment.GetPaymentsHandler)
	mux.HandleFunc("/api/payments/create", payment.CreatePaymentHandler)
	mux.HandleFunc("/api/payments/update", payment.UpdatePaymentHandler)
	mux.HandleFunc("/api/payments/delete", payment.DeletePaymentHandler)
	mux.HandleFunc("/api/payments/search", payment.SearchPaymentsHandler)

	handler := middleware.ChainMiddleware(
		mux,
		middleware.GzipMiddleware,
		middleware.CacheMiddleware(middleware.NewCacheConfig()),
		middleware.ApiLogMiddleware,
		middleware.TracingMiddleware,
		middleware.JWTMiddleware([]string{"/api/login", "/api/register", "/api/logout"}),
		middleware.CircuitBreakerMiddleware(10*time.Second),
		middleware.RateLimitMiddleware(1, 10),
		middleware.RequestContextMiddleware,
		middleware.CorsMiddleware,
	)

	// Create server with configured handler
	server := &http.Server{
		Addr:    fmt.Sprintf(":%v", cfg.GraphQLPort),
		Handler: handler,
		// Add timeouts to prevent slow clients from holding resources
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   15 * time.Second,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB

	}

	// Start server in a goroutine
	go func() {
		fmt.Printf("Server is running at http://localhost:%v\n", cfg.GraphQLPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v\n", err)
		}
	}()

	// Wait for interrupt signal
	<-ctx.Done()
	fmt.Println("\nShutting down server...")

	// Create shutdown context with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server forced to shutdown: %v\n", err)
	}

	fmt.Println("Server stopped gracefully")
}
