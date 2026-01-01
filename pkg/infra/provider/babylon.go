package provider

import (
	"fmt"
	"log"
	"time"

	"github.com/PhilippReinke/kino-berlin/pkg/domain"
	"github.com/gocolly/colly/v2"
)

type Babylon struct {
	c       *colly.Collector
	baseURL string
}

var _ domain.Provider = &Babylon{}

func NewBabylon() *Babylon {
	return &Babylon{
		c:       colly.NewCollector(),
		baseURL: "https://babylonberlin.eu",
	}
}

func (b Babylon) Name() string {
	return "Kino Babylon"
}

func (b Babylon) Scrape() ([]domain.Screening, error) {
	var screenings []domain.Screening

	b.c.OnHTML("#regridart-207", func(e *colly.HTMLElement) {
		e.ForEach("li", func(n int, e *colly.HTMLElement) {
			titles := e.ChildTexts("h3")
			if len(titles) <= 2 {
				return
			}

			date, err := parseDate(e.Attr("data-date"))
			if err != nil {
				log.Printf("Failed to parse date: %v", err)
				return
			}

			var duration time.Duration
			runtimeTexts := e.ChildTexts(".runtime")
			if len(runtimeTexts) > 0 {
				var err error
				duration, err = parseDuration(runtimeTexts[0])
				if err != nil {
					log.Printf("Failed to parse duration: %v", err)
				}
			}

			link := b.baseURL + e.ChildAttr(".mix-title", "href")

			title := e.ChildTexts("h3")[2]

			language := "" // TODO

			screeningID := domain.NewScreeningID(
				title,
				date,
				b.Name(),
				language,
			)

			screenings = append(screenings, domain.Screening{
				ID:          screeningID,
				Title:       title,
				Description: "",
				Start:       date,
				Duration:    duration,
				Cinema:      b.Name(),
				Language:    language,
				Links: domain.ScreeningLinks{
					Details:       link,
					ThumbnailLink: e.ChildAttr(".fancybox", "href"),
				},
				UpdatedAt: time.Now(),
			})
		})
	})

	if err := b.c.Visit(b.baseURL + "/programm"); err != nil {
		return []domain.Screening{}, fmt.Errorf("running colly: %w", err)
	}

	return screenings, nil
}

func parseDate(dateString string) (time.Time, error) {
	tz, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		return time.Time{}, fmt.Errorf("creating timezone: %w", err)
	}

	date, err := time.ParseInLocation("2006-01-02 15:04:05", dateString, tz)
	if err != nil {
		return time.Time{}, fmt.Errorf("parsing date %q: %w", dateString, err)
	}

	return date, nil
}

func parseDuration(durationString string) (time.Duration, error) {
	var minutes int

	_, err := fmt.Sscanf(durationString, "%d min.", &minutes)
	if err != nil {
		return 0, fmt.Errorf("parsing duration %q: %w", durationString, err)
	}

	return time.Duration(minutes) * time.Minute, nil
}
