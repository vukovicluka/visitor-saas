package dashboard

import (
	"encoding/json"
	"log"
	"net/http"
)

type Handler struct {
	queries *Queries
}

func NewHandler(queries *Queries) *Handler {
	return &Handler{queries: queries}
}

func (h *Handler) HandleSummary(w http.ResponseWriter, r *http.Request) {
	domain, days := parseParams(r)
	if domain == "" {
		http.Error(w, "domain is required", http.StatusBadRequest)
		return
	}

	stats, err := h.queries.Summary(r.Context(), domain, days)
	if err != nil {
		log.Printf("summary query error: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, stats)
}

func (h *Handler) HandlePages(w http.ResponseWriter, r *http.Request) {
	domain, days := parseParams(r)
	if domain == "" {
		http.Error(w, "domain is required", http.StatusBadRequest)
		return
	}

	pages, err := h.queries.Pages(r.Context(), domain, days)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, pages)
}

func (h *Handler) HandleReferrers(w http.ResponseWriter, r *http.Request) {
	domain, days := parseParams(r)
	if domain == "" {
		http.Error(w, "domain is required", http.StatusBadRequest)
		return
	}

	refs, err := h.queries.Referrers(r.Context(), domain, days)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, refs)
}

func parseParams(r *http.Request) (string, int) {
	domain := r.URL.Query().Get("domain")
	period := r.URL.Query().Get("period")

	days := 30
	switch period {
	case "today":
		days = 1
	case "7d":
		days = 7
	case "30d":
		days = 30
	case "12m":
		days = 365
	}

	return domain, days
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}
