package models

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"gorm.io/gorm"
)

var reAlphanum = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_-]{1,18}[a-zA-Z0-9]$`)

// Entry represents a log line in the database.
type Entry struct {
	gorm.Model
	Timestamp   time.Time
	LogLine     string
	Application string
	Host        string
	Environment string
	IP          string
}

// PrettyTimestamp prints entry timestamps in an easily readable format for the
// web GUI.
func (e Entry) PrettyTimestamp() string {
	return e.Timestamp.Format("2006-01-02T15:04:05Z07")
}

// EntryService defines what this package can provide to the user, namely
// interacting with Entries in the database.
type EntryService interface {
	// Create creates one or more entries in the database.
	Create(...Entry) error

	// Entries queries the database for entries given an EntriesQuery object,
	// which can help narrow down the search to certain parameters.
	Entries(EntriesQuery) ([]Entry, error)

	// Applications, Hosts and Environments fetches every unique application,
	// host and environment name from the entries table.
	Applications() ([]string, error)
	Hosts() ([]string, error)
	Environments() ([]string, error)

	// AutoMigrate prepares the database for storing Entries.
	AutoMigrate() error
	// DestructiveReset destroys all entry data in the database and
	// subsequently performs an AutoMigrate.
	DestructiveReset() error
}

// NewEntryService creates an EntryService given a database connection.
func NewEntryService(db *gorm.DB) EntryService {
	return &entryService{db}
}

type entryService struct {
	db *gorm.DB
}

// EntriesQuery has all the fields available for querying the database for
// specific categories of entries.
type EntriesQuery struct {
	Application string
	Host        string
	Environment string
	Page        int
	FromTime    time.Time
	ToTime      time.Time
}

func (es *entryService) Entries(q EntriesQuery) ([]Entry, error) {
	const pageSize = 100
	var ex []Entry
	if q.Page == 0 {
		q.Page = 1
	}
	db := es.db.Offset((q.Page - 1) * pageSize).Limit(pageSize)
	if !q.FromTime.IsZero() || !q.ToTime.IsZero() {
		db = db.Where("timestamp BETWEEN ? AND ?", q.FromTime, q.ToTime)
	}
	where := Entry{
		Application: q.Application,
		Host:        q.Host,
		Environment: q.Environment,
	}
	res := db.Where(where).Order("timestamp").Find(&ex)
	if err := res.Error; err != nil {
		return nil, err
	}
	return ex, nil
}

func (es *entryService) Applications() ([]string, error) {
	var ex []Entry
	err := es.db.Distinct("application").Order("application").Find(&ex).Error
	if err != nil {
		return nil, err
	}
	var res []string
	for _, e := range ex {
		res = append(res, e.Application)
	}
	return res, nil
}

func (es *entryService) Hosts() ([]string, error) {
	var ex []Entry
	err := es.db.Distinct("host").Order("host").Find(&ex).Error
	if err != nil {
		return nil, err
	}
	var res []string
	for _, e := range ex {
		res = append(res, e.Host)
	}
	return res, nil
}

func (es *entryService) Environments() ([]string, error) {
	var ex []Entry
	err := es.db.Distinct("environment").Order("environment").Find(&ex).Error
	if err != nil {
		return nil, err
	}
	var res []string
	for _, e := range ex {
		res = append(res, e.Environment)
	}
	return res, nil
}

func (es *entryService) Create(ex ...Entry) error {
	if err := validateEntries(ex...); err != nil {
		return fmt.Errorf("validating entries: %w", err)
	}
	return es.db.Create(&ex).Error
}

func validateEntries(ex ...Entry) error {
	for _, e := range ex {
		if e.LogLine == "" {
			return errors.New("field LogLine must not be empty")
		}
		if e.Timestamp.IsZero() {
			return errors.New("field Timestamp must not be empty")
		}
		if !reAlphanum.MatchString(e.Application) ||
			!reAlphanum.MatchString(e.Environment) ||
			!reAlphanum.MatchString(e.Host) {
			return fmt.Errorf("application, host and environment must match %s", reAlphanum)
		}
	}
	return nil
}

func (es *entryService) AutoMigrate() error {
	return es.db.AutoMigrate(&Entry{})
}

func (es *entryService) DestructiveReset() error {
	if err := es.db.Migrator().DropTable(&Entry{}); err != nil {
		return err
	}
	return es.AutoMigrate()
}
