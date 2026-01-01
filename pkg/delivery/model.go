package delivery

import "time"

type ScreeningViewModel struct {
	Title         string
	Cinema        string
	Duration      int
	Date          time.Time
	Link          string
	ThumbnailLink string
}
