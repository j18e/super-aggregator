package models

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Entry struct {
	gorm.Model
	Timestamp   time.Time
	LogLine     string
	Application string
	Host        string
	Environment string
	IP          string
}

func (e Entry) PrettyTimestamp() string {
	return e.Timestamp.Format("2006-01-02T15:04:05Z07")
}

type EntryService interface {
	Create([]Entry) error
	Entries(EntriesQuery) ([]Entry, error)
	Applications() ([]string, error)
	Hosts() ([]string, error)
	Environments() ([]string, error)

	AutoMigrate() error
	DestructiveReset() error
}

func NewEntryService(db *gorm.DB) EntryService {
	return &entryService{db}
}

type entryService struct {
	db *gorm.DB
}

type EntriesQuery struct {
	Application string
	Host        string
	Environment string
	Page        int
}

func (es *entryService) Entries(q EntriesQuery) ([]Entry, error) {
	const pageSize = 100
	var ex []Entry
	if q.Page == 0 {
		q.Page = 1
	}
	db := es.db.Offset((q.Page - 1) * pageSize).Limit(pageSize)
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

func (es *entryService) Create(ex []Entry) error {
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
