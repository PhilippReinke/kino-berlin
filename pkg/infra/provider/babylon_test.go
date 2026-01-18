package provider

import (
	"testing"
)

func TestBabylon_Name(t *testing.T) {
	b := NewBabylon()
	want := "Kino Babylon"
	if got := b.Name(); got != want {
		t.Errorf("Name() = %q, want %q", got, want)
	}
}

func TestBabylon_Scrape(t *testing.T) {
	// TODO
}
