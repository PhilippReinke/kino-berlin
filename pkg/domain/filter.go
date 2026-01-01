package domain

import "time"

// Filter returns true if filter matches.
//
// As an example a date filter would return true if a screening matches the
// date.
type Filter func(Screening) bool

func DateFilter(date time.Time) Filter {
	year, month, day := date.Date()
	return func(s Screening) bool {
		sYear, sMonth, sDay := s.Start.Date()
		return sYear == year && sMonth == month && sDay == day
	}
}

func ExpiredFilter(maxAge time.Duration) Filter {
	return func(s Screening) bool {
		return time.Since(s.UpdatedAt) < maxAge
	}
}

func ExpiredScreeningFilter() Filter {
	return func(s Screening) bool {
		return s.Start.After(time.Now())
	}
}

func CinemaFilter(cinema string) Filter {
	return func(s Screening) bool {
		return s.Cinema == cinema
	}
}
