package handlers

import (
	"encoding/json"
	"go-kasir-api/services"
	"net/http"
)

type ReportHandler struct {
	service *services.ReportService
}

func NewReportHandler(service *services.ReportService) *ReportHandler {
	return &ReportHandler{service: service}
}

// GET api/report?start_date=...&end_date=...
func (h *ReportHandler) HandleReport(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetReport(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *ReportHandler) GetReport(w http.ResponseWriter, r *http.Request) {
	start := r.URL.Query().Get("start_date")
	end := r.URL.Query().Get("end_date")

	if start == "" || end == "" {
		http.Error(w, "start_date and end_date are required", http.StatusBadRequest)
		return
	}

	report, err := h.service.GetReport(start, end)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}
