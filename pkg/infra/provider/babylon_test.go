package provider

import (
	"testing"
)

func TestBabylon_Name(t *testing.T) {
	b := NewBabylon()
	if got := b.Name(); got != "Kino Babylon" {
		t.Errorf("Name() = %q, want %q", got, "Kino Babylon")
	}
}

func TestBabylon_Scrape(t *testing.T) {
	// TODO
}
