// cmd/server/main.go
package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/go-openapi/runtime/middleware"
	"movie-project/config"
	"movie-project/internal/handler"
	"movie-project/internal/repository"
	"movie-project/internal/service"
	"movie-project/pkg/logger"
	"movie-project/pkg/metrics"
	pb "movie-project/proto/movie"
)

func main() {
	// Initialize logger
	log := logger.NewLogger()

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Error("Failed to load config", "error", err)
		os.Exit(1)
	}

	// Initialize database connection
	db, err := gorm.Open(postgres.Open(cfg.GetDSN()), &gorm.Config{})
	if err != nil {
		log.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}

	// Initialize repository, service, and handler
	repo := repository.NewMovieRepository(db, log)
	svc := service.NewMovieService(repo, log)
	movieHandler := handler.NewMovieHandler(svc, log)

	// Initialize Prometheus metrics
	metrics.InitMetrics()

	// Initialize gRPC server
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(metrics.UnaryServerInterceptor),
	)
	pb.RegisterMovieServiceServer(grpcServer, movieHandler)
	reflection.Register(grpcServer)

	// Start gRPC server
	grpcAddr := fmt.Sprintf("%s:%s", cfg.ServerHost, cfg.GRPCPort)
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Error("Failed to listen for gRPC", "error", err)
		os.Exit(1)
	}
	go func() {
		log.Info("Starting gRPC server", "address", grpcAddr)
		if err := grpcServer.Serve(lis); err != nil {
			log.Error("Failed to serve gRPC", "error", err)
			os.Exit(1)
		}
	}()

	// Initialize gRPC-Gateway
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	gwmux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err = pb.RegisterMovieServiceHandlerFromEndpoint(ctx, gwmux, grpcAddr, opts)
	if err != nil {
		log.Error("Failed to register gRPC-Gateway", "error", err)
		os.Exit(1)
	}

	// Create an HTTP server
	mux := http.NewServeMux()
	mux.Handle("/", instrumentHandler(gwmux, "grpc_gateway"))
	mux.Handle("/metrics", promhttp.Handler())

	// Configure CORS
	corsMiddleware := corsMiddleware(cfg.AllowedOrigins)

	// Настройка Swagger UI
	httpAddr := fmt.Sprintf("%s:%s", cfg.ServerHost, cfg.ServerPort)
	swaggerOpts := middleware.SwaggerUIOpts{SpecURL: fmt.Sprintf("http://%s/api/api.swagger.json", httpAddr)}
	sh := middleware.SwaggerUI(swaggerOpts, nil)
	mux.Handle("/docs", sh)
	mux.HandleFunc("/api/api.swagger.json", func(w http.ResponseWriter, r *http.Request) {
		log.Info("Received request for Swagger JSON", "path", r.URL.Path)
		http.ServeFile(w, r, "api/api.swagger.json")
	})

	srv := &http.Server{
		Addr:    httpAddr,
		Handler: corsMiddleware(mux),
	}

	// Start HTTP server
	go func() {
		log.Info("Starting HTTP server", "address", httpAddr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("Failed to serve HTTP", "error", err)
			os.Exit(1)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down server...")

	// Shutdown HTTP server
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Error("Server forced to shutdown", "error", err)
	}

	// Stop gRPC server
	grpcServer.GracefulStop()

	log.Info("Server exited")
}

func corsMiddleware(allowedOrigins []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*") // Разрешаем все источники
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Max-Age", "3600")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func instrumentHandler(next http.Handler, handlerName string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := newResponseWriter(w)
		next.ServeHTTP(rw, r)

		duration := time.Since(start).Seconds()
		statusCode := rw.statusCode

		metrics.RequestDuration.WithLabelValues(handlerName, r.Method).Observe(duration)
		metrics.TotalRequests.WithLabelValues(handlerName, strconv.Itoa(statusCode), r.Method).Inc()
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
