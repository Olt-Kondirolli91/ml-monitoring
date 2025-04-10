package server

import (
    "context"
    "database/sql"
    "log"
    "net/http"

    "github.com/Olt-Kondirolli91/ml-monitoring/internal/repository"
    "github.com/gorilla/mux"
)

// Server holds references to repositories and the router
type Server struct {
    InferenceRepo repository.InferenceRepository
    FeedbackRepo  repository.FeedbackRepository
    Router        *mux.Router
    httpServer    *http.Server
}

// NewServer creates a new Server instance with the given repositories
func NewServer(db *sql.DB) *Server {
    infRepo := repository.NewInferenceRepository(db)
    fbRepo := repository.NewFeedbackRepository(db)

    s := &Server{
        InferenceRepo: infRepo,
        FeedbackRepo:  fbRepo,
        Router:        mux.NewRouter(),
    }
    s.routes()
    return s
}

// routes sets up our HTTP endpoints
func (s *Server) routes() {
    // Health check
    s.Router.HandleFunc("/health", s.handleHealth).Methods("GET")

    // Inference endpoints
    s.Router.HandleFunc("/inferences", s.handleCreateInference).Methods("POST")
    s.Router.HandleFunc("/inferences/{id}", s.handleGetInference).Methods("GET")

    // Feedback endpoint
    s.Router.HandleFunc("/inferences/{id}/feedback", s.handleCreateFeedback).Methods("POST")
    s.Router.HandleFunc("/inferences/{id}/feedback", s.handleGetFeedback).Methods("GET")
}

// starts the HTTP server on the specified port
func (s *Server) Start(port string) {
    s.httpServer = &http.Server{
        Addr:    ":" + port,
        Handler: s.Router,
    }

    log.Printf("Starting HTTP server on port %s\n", port)
    if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        log.Fatalf("Could not listen on port %s: %v\n", port, err)
    }
}

// shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
    log.Println("Shutting down server...")
    return s.httpServer.Shutdown(ctx)
}
