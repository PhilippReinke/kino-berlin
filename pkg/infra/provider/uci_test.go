package provider

import (
	"testing"
)

func TestUCI_Name(t *testing.T) {
	u := NewUCI()
	want := "UCI Berlin"
	if got := u.Name(); got != want {
		t.Errorf("Name() = %q, want %q", got, want)
	}
}

func TestUCI_Scrape(t *testing.T) {
	// TODO
}

func TestUCI_Manual(t *testing.T) {
	// t.Skip("manual test")

	// Usage:
	// go test -v -run TestUCI_Manual pkg/infra/provider/*.go
	u := NewUCI()

	screenings, err := u.Scrape()
	if err != nil {
		t.Errorf("Scrape() = %v", err)
	}

	for _, s := range screenings {
		t.Logf("%+v\n", s)
	}
}
