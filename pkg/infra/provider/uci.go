package provider

import (
	"fmt"

	"github.com/PhilippReinke/kino-berlin/pkg/domain"
	"github.com/gocolly/colly/v2"
)

/*

current links:

https://www.uci-kinowelt.de/kinoprogramm

- https://www.uci-kinowelt.de/kinoprogramm/berlin-gropius-passagen/43/list
- https://www.uci-kinowelt.de/kinoprogramm/berlin-east-side-gallery/82/list
- https://www.uci-kinowelt.de/kinoprogramm/berlin-am-eastgate/44/list
- https://www.uci-kinowelt.de/kinoprogramm/potsdam/59/list

start with div:

	class="film show  "

*/

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

	u.c.OnHTML("#scheduleContainer", func(e *colly.HTMLElement) {
		e.ForEach("div.film.show", func(n int, el *colly.HTMLElement) {
			title := el.ChildText(".eventkalender--item--description--text--eventtitle a")

			screenings = append(screenings, domain.Screening{
				Title: title,
			})
		})
	})

	if err := u.c.Visit(u.baseURL + "/kinoprogramm/berlin-east-side-gallery/82/list"); err != nil {
		return []domain.Screening{}, fmt.Errorf("running colly: %w", err)
	}

	return screenings, nil
}
