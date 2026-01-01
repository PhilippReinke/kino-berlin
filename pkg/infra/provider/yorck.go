package provider

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PhilippReinke/kino-berlin/pkg/domain"
	"github.com/PhilippReinke/kino-berlin/pkg/infra/provider/yorckmodel"
	"github.com/gocolly/colly/v2"
)

const (
	scriptTagBegin = "<script id=\"__NEXT_DATA__\" type=\"application/json\">"
	scriptTagEnd   = "</script>"
)

type Yorck struct {
	c       *colly.Collector
	baseURL string
}

var _ domain.Provider = &Yorck{}

func NewYorck() *Yorck {
	return &Yorck{
		c:       colly.NewCollector(),
		baseURL: "https://www.yorck.de",
	}
}

func (y Yorck) Name() string {
	return "Yorck Kinos"
}

func (y Yorck) Scrape() ([]domain.Screening, error) {
	yorckAddress := fmt.Sprintf("%v/%v", y.baseURL, "filme")

	res, err := http.Get(yorckAddress)
	if err != nil {
		return []domain.Screening{}, fmt.Errorf("fetching from %q: %w", yorckAddress, err)
	}
	defer res.Body.Close()

	bodyByte, err := io.ReadAll(res.Body)
	if err != nil {
		return []domain.Screening{}, fmt.Errorf("reading body: %w", err)
	}
	body := string(bodyByte)

	begin := strings.Index(body, scriptTagBegin)
	if begin == -1 {
		return []domain.Screening{}, fmt.Errorf("finding begin of film data")
	}
	begin += len(scriptTagBegin)
	end := strings.Index(body[begin:], scriptTagEnd)
	if end == -1 {
		return []domain.Screening{}, fmt.Errorf("finding end of film data")
	}

	jsonString := body[begin : begin+end]

	var films yorckmodel.FilmsYorck
	if err := json.Unmarshal([]byte(jsonString), &films); err != nil {
		return []domain.Screening{}, fmt.Errorf("unmarshaling JSON: %w", err)
	}

	tz, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		return []domain.Screening{}, fmt.Errorf("creating timezone: %w", err)
	}

	var screenings []domain.Screening
	for _, film := range films.Props.PageProps.Films {
		for _, session := range film.Fields.Sessions {
			title := film.Fields.Title
			start := time.Date(
				session.Fields.StartTime.Year(),
				session.Fields.StartTime.Month(),
				session.Fields.StartTime.Day(),
				session.Fields.StartTime.Hour(),
				session.Fields.StartTime.Minute(),
				session.Fields.StartTime.Second(),
				session.Fields.StartTime.Nanosecond(),
				tz,
			)
			cinema := session.Fields.Cinema.Fields.Name
			language := ""
			screeningID := domain.NewScreeningID(
				title,
				start,
				cinema,
				language,
			)
			duration := time.Minute * time.Duration(film.Fields.Runtime)

			screenings = append(screenings, domain.Screening{
				ID:          screeningID,
				Title:       title,
				Description: "",
				Start:       start,
				Duration:    duration,
				Cinema:      cinema,
				Language:    language,
				Links: domain.ScreeningLinks{
					Details:       fmt.Sprintf("%v/%v", yorckAddress, film.Fields.Slug),
					ThumbnailLink: createThumbnailLink(film.Fields.HeroImage.Fields.Image.FieldsImage.File.URL),
				},
				UpdatedAt: time.Now(),
			})
		}
	}

	return screenings, nil
}

func createThumbnailLink(thumbnailURL string) string {
	u, _ := url.Parse("https:" + thumbnailURL)
	q := u.Query()
	q.Set("w", "480")
	q.Set("q", "75")
	u.RawQuery = q.Encode()

	return u.String()
}
