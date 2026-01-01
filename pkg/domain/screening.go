package domain

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

type Screening struct {
	ID          ScreeningID
	Title       string
	Description string
	Start       time.Time
	Duration    time.Duration
	Cinema      string
	Language    string
	Links       ScreeningLinks
	UpdatedAt   time.Time
}

type ScreeningLinks struct {
	Details       string
	ThumbnailLink string
}

type ScreeningID string

func NewScreeningID(
	title string,
	start time.Time,
	cinema string,
	language string,
) ScreeningID {
	startString := start.Format(time.RFC3339)
	hashData := title + startString + cinema + language

	sum := sha256.Sum256([]byte(hashData))
	return ScreeningID(hex.EncodeToString(sum[:]))
}
