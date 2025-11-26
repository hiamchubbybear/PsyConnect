package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/loggingservice/pkg/settings"
	"go.uber.org/zap"
)

// Server represents the HTTP server for health checks and metrics
type Server struct {
	config        settings.ServerConfig
	logger        *zap.Logger
	httpServer    *http.Server
	messagesCount int64
	errorsCount   int64
	startTime     time.Time
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status         string    `json:"status"`
	Uptime         string    `json:"uptime"`
	MessagesCount  int64     `json:"messages_count"`
	ErrorsCount    int64     `json:"errors_count"`
	Timestamp      time.Time `json:"timestamp"`
	Version        string    `json:"version"`
}

// MetricsResponse represents the metrics response
type MetricsResponse struct {
	MessagesProcessed int64   `json:"messages_processed"`
	ErrorsCount       int64   `json:"errors_count"`
	ErrorRate         float64 `json:"error_rate"`
	Uptime            string  `json:"uptime"`
}

// NewServer creates a new HTTP server
func NewServer(config settings.ServerConfig, logger *zap.Logger) *Server {
	return &Server{
		config:    config,
		logger:    logger,
		startTime: time.Now(),
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", s.handleHealth)
	mux.HandleFunc("/metrics", s.handleMetrics)
	mux.HandleFunc("/ready", s.handleReady)

	s.httpServer = &http.Server{
		Addr:         fmt.Sprintf(":%d", s.config.Port),
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	s.logger.Info("Starting HTTP server", zap.Int("port", s.config.Port))

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Error("HTTP server error", zap.Error(err))
		}
	}()

	return nil
}

// Stop stops the HTTP server
func (s *Server) Stop() error {
	if s.httpServer != nil {
		s.logger.Info("Stopping HTTP server")
		return s.httpServer.Close()
	}
	return nil
}

// IncrementMessages increments the messages counter
func (s *Server) IncrementMessages() {
	atomic.AddInt64(&s.messagesCount, 1)
}

// IncrementErrors increments the errors counter
func (s *Server) IncrementErrors() {
	atomic.AddInt64(&s.errorsCount, 1)
}

// handleHealth handles health check requests
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	uptime := time.Since(s.startTime)

	response := HealthResponse{
		Status:        "healthy",
		Uptime:        uptime.String(),
		MessagesCount: atomic.LoadInt64(&s.messagesCount),
		ErrorsCount:   atomic.LoadInt64(&s.errorsCount),
		Timestamp:     time.Now(),
		Version:       "1.0.0",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// handleMetrics handles metrics requests
func (s *Server) handleMetrics(w http.ResponseWriter, r *http.Request) {
	uptime := time.Since(s.startTime)
	messages := atomic.LoadInt64(&s.messagesCount)
	errors := atomic.LoadInt64(&s.errorsCount)

	var errorRate float64
	if messages > 0 {
		errorRate = float64(errors) / float64(messages) * 100
	}

	response := MetricsResponse{
		MessagesProcessed: messages,
		ErrorsCount:       errors,
		ErrorRate:         errorRate,
		Uptime:            uptime.String(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// handleReady handles readiness probe requests
func (s *Server) handleReady(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ready"})
}
