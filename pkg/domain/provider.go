package domain

type Provider interface {
	Scrape() ([]Screening, error)
	Name() string
}
