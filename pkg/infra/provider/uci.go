package provider

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PhilippReinke/kino-berlin/pkg/domain"
	"github.com/gocolly/colly/v2"
)

var uciCinemas = []struct {
	path string
	name string
}{
	{"/kinoprogramm/berlin-gropius-passagen/43/list", "UCI Gropius Passagen"},
	{"/kinoprogramm/berlin-east-side-gallery/82/list", "UCI East Side Gallery"},
	{"/kinoprogramm/berlin-am-eastgate/44/list", "UCI Am Eastgate"},
	{"/kinoprogramm/potsdam/59/list", "UCI Potsdam"},
}

type UCI struct {
	c       *colly.Collector
	baseURL string
}

var _ domain.Provider = &UCI{}

func NewUCI() *UCI {
	return &UCI{
		c:       colly.NewCollector(),
		baseURL: "https://www.uci-kinowelt.de",
	}
}

func (u UCI) Name() string {
	return "UCI Berlin"
}

func (u UCI) Scrape() ([]domain.Screening, error) {
	var screenings []domain.Screening
	var currentCinema string

	u.c.OnHTML("#scheduleContainer", func(e *colly.HTMLElement) {
		e.ForEach("div.film.show", func(_ int, el *colly.HTMLElement) {
			title := el.ChildText("h2.eventkalender--item--description--text--eventtitle a")
			detailsHref := el.ChildAttr("h2.eventkalender--item--description--text--eventtitle a", "href")
			detailsLink := u.baseURL + detailsHref
			thumbnail := el.ChildAttr(".eventkalender--item--description--picture img", "src")
			duration := uciParseDuration(el)

			el.ForEach("a.performance", func(_ int, perf *colly.HTMLElement) {
				dateStr := perf.Attr("data-date")
				timeStr := strings.Trim(perf.Attr("data-time"), "'")
				if dateStr == "" || timeStr == "" {
					return
				}
				start, err := parseUCIDateTime(dateStr, timeStr)
				if err != nil {
					log.Printf("Failed to parse UCI date/time %q %q: %v", dateStr, timeStr, err)
					return
				}

				language := ""
				if strings.Contains(perf.Attr("class"), "attribut-ov") {
					language = "OV"
				}

				screenings = append(screenings, domain.Screening{
					ID:          domain.NewScreeningID(title, start, currentCinema, language),
					Title:       title,
					Description: "",
					Start:       start,
					Duration:    duration,
					Cinema:      currentCinema,
					Language:    language,
					Links: domain.ScreeningLinks{
						Details:       detailsLink,
						ThumbnailLink: thumbnail,
					},
					UpdatedAt: time.Now(),
				})
			})
		})
	})

	for _, c := range uciCinemas {
		currentCinema = c.name
		if err := u.c.Visit(u.baseURL + c.path); err != nil {
			return nil, fmt.Errorf("visiting %q: %w", c.path, err)
		}
	}

	return screenings, nil
}

func parseUCIDateTime(dateStr, timeStr string) (time.Time, error) {
	tz, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		return time.Time{}, fmt.Errorf("creating timezone: %w", err)
	}
	return time.ParseInLocation("20060102 15:04", dateStr+" "+timeStr, tz)
}

func uciParseDuration(el *colly.HTMLElement) time.Duration {
	for _, s := range el.ChildTexts("ul.film-info.infolist li") {
		var minutes int
		if _, err := fmt.Sscanf(s, "%dmin", &minutes); err == nil {
			return time.Duration(minutes) * time.Minute
		}
		if _, err := fmt.Sscanf(s, "%d min.", &minutes); err == nil {
			return time.Duration(minutes) * time.Minute
		}
	}
	return 0
}
