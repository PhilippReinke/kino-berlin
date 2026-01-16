package delivery

import (
	"log"
	"net/http"
	"time"

	"github.com/PhilippReinke/kino-berlin/pkg/domain"
)

var (
	defaultFilters = []domain.Filter{
		domain.RecentlyUpdated(47 * time.Hour),
		domain.AlreadyOver(),
	}
)

func (h *Handler) handleSelects(w http.ResponseWriter, r *http.Request) {
	cinemas, err := h.app.GetAvailableCinemas(defaultFilters...)
	if err != nil {
		h.renderError(w, err)
		return
	}

	dates, err := h.app.GetAvailableDates(defaultFilters...)
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

	filters := defaultFilters

	if dateStr := r.FormValue("dates"); dateStr != "" {
		if date, err := time.Parse(time.DateOnly, dateStr); err == nil {
			filters = append(filters, domain.DateMatches(date))
		}
	}

	if cinema := r.FormValue("cinemas"); cinema != "" {
		filters = append(filters, domain.CinemaMatches(cinema))
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
