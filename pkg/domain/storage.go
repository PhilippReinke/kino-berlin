package domain

type Storage interface {
	Upsert(screenings Screening) error
	Fetch(filter ...Filter) ([]Screening, error)
}
