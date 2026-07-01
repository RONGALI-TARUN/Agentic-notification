package server

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/example/notification-service/pkg/config"
	"github.com/example/notification-service/pkg/notification"
	"github.com/example/notification-service/pkg/notification/impl"
	"github.com/example/notification-service/pkg/util"
)

type Server struct {
	cfg     *config.Config
	service notification.NotificationService
	http    *http.Server
}

func NewServer(cfg *config.Config) *Server {
	return &Server{
		cfg:     cfg,
		service: impl.NewNotificationService(cfg),
	}
}

func (s *Server) Run() error {
	mux := http.NewServeMux()
	mux.Handle("/notify", s.loggingMiddleware(http.HandlerFunc(s.notifyHandler)))
	mux.Handle("/health", s.loggingMiddleware(http.HandlerFunc(s.healthHandler)))
	mux.Handle("/metrics", s.loggingMiddleware(http.HandlerFunc(s.metricsHandler)))

	s.http = &http.Server{
		Addr:    s.cfg.Address,
		Handler: mux,
	}

	log.Printf("notification service listening on %s", s.cfg.Address)
	return s.http.ListenAndServe()
}

func (s *Server) notifyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		util.WriteJSON(w, http.StatusMethodNotAllowed, util.ErrorResponse{Error: "method not allowed"})
		return
	}

	var req notification.NotificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.WriteJSON(w, http.StatusBadRequest, util.ErrorResponse{Error: "invalid request payload"})
		return
	}

	if err := s.service.Send(req); err != nil {
		log.Printf("notification error: %v", err)
		util.WriteJSON(w, http.StatusInternalServerError, util.ErrorResponse{Error: err.Error()})
		return
	}

	util.WriteJSON(w, http.StatusAccepted, map[string]string{"status": "accepted"})
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	util.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) metricsHandler(w http.ResponseWriter, r *http.Request) {
	util.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"requests":      0,
		"notifications": 0,
	})
}

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("incoming request method=%s path=%s remote=%s", r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(w, r)
		log.Printf("completed request method=%s path=%s duration=%s", r.Method, r.URL.Path, time.Since(start))
	})
}
