package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"4imdb-seasons-tracker/internal/service"
)

type Handler struct {
	service *service.TrackerService
	logger  *log.Logger
}

func NewHandler(service *service.TrackerService, logger *log.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/health", h.healthCheck)
	mux.HandleFunc("/api/v1/series", h.handleSeries)
	mux.HandleFunc("/api/v1/check", h.triggerCheck)
}

func (h *Handler) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *Handler) handleSeries(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.listSeries(w, r)
	case http.MethodPost:
		h.addSeries(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

type AddSeriesRequest struct {
	ImdbID string `json:"imdb_id"`
	Season int    `json:"season"`
}

func (h *Handler) addSeries(w http.ResponseWriter, r *http.Request) {
	var req AddSeriesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.ImdbID == "" || req.Season <= 0 {
		http.Error(w, "Invalid imdb_id or season", http.StatusBadRequest)
		return
	}

	msg, err := h.service.AddSeries(req.ImdbID, req.Season)
	if err != nil {
		h.logger.Printf("Error adding series: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": msg})
}

func (h *Handler) listSeries(w http.ResponseWriter, r *http.Request) {
	h.logger.Printf("List series request from %s", r.URL.Path)

	series, err := h.service.GetAll()
	if err != nil {
		h.logger.Printf("Error getting series: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(series)
}

func (h *Handler) triggerCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	go func() {
		if err := h.service.CheckAll(); err != nil {
			h.logger.Printf("Error during check: %v", err)
		}
	}()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Check started"})
}
