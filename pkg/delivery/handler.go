package delivery

import (
	"log"
	"net/http"
	"time"

	"github.com/PhilippReinke/kino-berlin/pkg/domain"
)

func (h *Handler) handleSelects(w http.ResponseWriter, r *http.Request) {
	cinemas, err := h.app.GetAvailableCinemas()
	if err != nil {
		h.renderError(w, err)
		return
	}

	dates, err := h.app.GetAvailableDates()
	if err != nil {
		h.renderError(w, err)
		return
	}

	data := struct {
		ScrapeIDs []string
		Cinemas   []string
		Dates     []time.Time
	}{
		ScrapeIDs: []string{},
		Cinemas:   cinemas,
		Dates:     dates,
	}

	if err := h.templates.ExecuteTemplate(w, "selects", data); err != nil {
		h.renderError(w, err)
		return
	}
}

func (h *Handler) handleScreenings(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.renderError(w, err)
		return
	}

	filters := []domain.Filter{
		domain.ExpiredFilter(47 * time.Hour),
		domain.ExpiredScreeningFilter(),
	}

	if dateStr := r.FormValue("dates"); dateStr != "" {
		if date, err := time.Parse(time.DateOnly, dateStr); err == nil {
			filters = append(filters, domain.DateFilter(date))
		}
	}

	if cinema := r.FormValue("cinemas"); cinema != "" {
		filters = append(filters, domain.CinemaFilter(cinema))
	}

	screenings, err := h.app.FetchScreenings(filters...)
	if err != nil {
		h.renderError(w, err)
		return
	}

	viewModels := make([]ScreeningViewModel, len(screenings))
	for i, s := range screenings {
		viewModels[i] = ScreeningViewModel{
			Title:         s.Title,
			Cinema:        s.Cinema,
			Duration:      int(s.Duration.Minutes()),
			Date:          s.Start,
			Link:          s.Links.Details,
			ThumbnailLink: s.Links.ThumbnailLink,
		}
	}

	if err := h.templates.ExecuteTemplate(w, "screenings", viewModels); err != nil {
		h.renderError(w, err)
		return
	}
}

func (h *Handler) renderError(w http.ResponseWriter, err error) {
	log.Printf("Error: %v", err)
	if err := h.templates.ExecuteTemplate(w, "error", err.Error()); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
