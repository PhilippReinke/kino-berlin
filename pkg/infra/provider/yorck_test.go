package provider

import "testing"

func TestYorck_Name(t *testing.T) {
	b := NewBabylon()
	want := "Yorck Kinos"
	if got := b.Name(); got != want {
		t.Errorf("Name() = %q, want %q", got, want)
	}
}

func TestYorck_Scrape(t *testing.T) {
	// TODO
}
