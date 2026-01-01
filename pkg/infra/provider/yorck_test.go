package provider

import "testing"

func TestYorck_Name(t *testing.T) {
	b := NewBabylon()
	if got := b.Name(); got != "Yorck Kinos" {
		t.Errorf("Name() = %q, want %q", got, "Yorck Kinos")
	}
}

func TestYorck_Scrape(t *testing.T) {
	// TODO
}
