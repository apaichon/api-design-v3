package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
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
	// Parse command line arguments for number of instances
	var instances int
	flag.IntVar(&instances, "n", 1, "Number of server instances to run")
	flag.Parse()

	// If passed as argument without flag, check os.Args
	if len(os.Args) > 1 && instances == 1 {
		if n, err := strconv.Atoi(os.Args[1]); err == nil {
			instances = n
		}
	}

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

	// WaitGroup to wait for all servers to stop
	var wg sync.WaitGroup

	// Start multiple server instances
	for i := 0; i < instances; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			port := cfg.GraphQLPort + index
			runServer(ctx, port)
		}(i)
	}

	// Wait for interrupt signal
	<-ctx.Done()
	log.Println("\nShutting down servers...")

	// Wait for all servers to stop
	wg.Wait()
	log.Println("All servers stopped gracefully")
}

// runServer starts a single server instance on the specified port
func runServer(ctx context.Context, port int) {
	// Create a new server instance
	server := createServer(port)

	// Start server
	go func() {
		log.Printf("Server is running at http://localhost:%v\n", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Server on port %d failed: %v\n", port, err)
		}
	}()

	<-ctx.Done()
	shutdownServer(server)
}

// createServer creates and configures a new HTTP server
func createServer(port int) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})
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
		middleware.JWTMiddleware([]string{"/api/health", "/api/login", "/api/register", "/api/logout"}),
		middleware.CircuitBreakerMiddleware(10*time.Second),
		middleware.RateLimitMiddleware(1, 10),
		middleware.RequestContextMiddleware,
		middleware.CorsMiddleware,
	)

	return &http.Server{
		Addr:           fmt.Sprintf(":%v", port),
		Handler:        handler,
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   15 * time.Second,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}
}

// shutdownServer gracefully shuts down a server
func shutdownServer(server *http.Server) {
	// Create shutdown context with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server on %s forced to shutdown: %v\n", server.Addr, err)
	} else {
		log.Printf("Server on %s stopped gracefully\n", server.Addr)
	}
}
