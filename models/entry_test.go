package models

import (
	"testing"
	"time"
)

func Test_validateEntries(t *testing.T) {
	validEntry := func() Entry {
		return Entry{
			Timestamp:   time.Now(),
			LogLine:     "something happened",
			Application: "app1",
			Environment: "dev",
			Host:        "server-1",
		}
	}
	valid := []Entry{
		validEntry(),
	}
	for _, e := range valid {
		if err := validateEntries(e); err != nil {
			t.Errorf("expected %v to be valid but was not: %s", e, err)
		}
	}
	invalid := []Entry{
		{},
		{LogLine: "something happened"},
	}
	invalid = append(invalid, func() Entry {
		e := validEntry()
		e.Host = "pc"
		return e
	}())
	invalid = append(invalid, func() Entry {
		e := validEntry()
		e.Application = "myapp-"
		return e
	}())
	for _, e := range invalid {
		if err := validateEntries(e); err == nil {
			t.Errorf("expected %v not to be valid but it was: %s", e, err)
		}
	}
}
