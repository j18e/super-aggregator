package models

import (
	"reflect"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func validEntry() Entry {
	return Entry{
		Timestamp:   time.Now().In(time.UTC),
		LogLine:     "something happened",
		Application: "app1",
		Environment: "dev",
		Host:        "server-1",
	}
}

func resetModel(e *Entry) {
	e.ID = 0
	e.CreatedAt = time.Time{}
	e.UpdatedAt = time.Time{}
	e.DeletedAt = gorm.DeletedAt{}
}

func TestEntryService_Entries(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("test.db"))
	if err != nil {
		t.Fatalf("opening gorm connection: %s", err)
	}
	es := NewEntryService(db)
	if err := es.DestructiveReset(); err != nil {
		t.Fatalf("automigrating database: %s", err)
	}

	exp := validEntry()

	if err := es.Create(exp); err != nil {
		t.Fatalf("creating entry: %s", err)
	}

	got, err := es.Entries(EntriesQuery{})
	if err != nil {
		t.Fatalf("did not expect err: %s", err)
	}
	if len(got) != 1 {
		t.Fatalf("expected 1 result, got %d", len(got))
	}
	resetModel(&got[0])
	if !reflect.DeepEqual(exp, got[0]) {
		t.Errorf("expected:\n%v\ngot:\n%v", exp, got)
	}
}

func Test_validateEntries(t *testing.T) {
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
