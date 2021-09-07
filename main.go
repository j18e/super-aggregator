package main

import (
	"time"

	"github.com/j18e/super-aggregator/controllers"
	"github.com/j18e/super-aggregator/models"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	db, err := gorm.Open(sqlite.Open("tmp/data.sqlite"), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		return err
	}
	db.AutoMigrate(&models.Entry{})
	es := models.NewEntryService(db)

	// if err := createEntries(es); err != nil {
	// 	return err
	// }

	ec := controllers.NewEntryController(es)

	r := gin.Default()
	r.LoadHTMLGlob("./views/*.html")
	r.GET("/", ec.EntriesHandler())
	r.POST("/api/entry", ec.CreateHandler())
	return r.Run(":9000")
}

func createEntries(es models.EntryService) error {
	ts := time.Now()
	for i := 0; i < 100; i++ {
		if err := es.Create(models.Entry{
			Timestamp:   ts,
			LogLine:     "something exciting happened",
			Application: "app1",
			Host:        "host1",
		}); err != nil {
			return err
		}
		ts = ts.Add(time.Hour * -6)
		if err := es.Create(models.Entry{
			Timestamp:   ts,
			LogLine:     "something boring happened",
			Application: "app2",
			Host:        "host2",
		}); err != nil {
			return err
		}
		ts = ts.Add(time.Hour * -6)
		if err := es.Create(models.Entry{
			Timestamp:   ts,
			LogLine:     "something whatever happened",
			Application: "app3",
			Host:        "host1",
		}); err != nil {
			return err
		}
		ts = ts.Add(time.Hour * -6)
		if err := es.Create(models.Entry{
			Timestamp:   ts,
			LogLine:     "nothing happened",
			Application: "app4",
			Host:        "host2",
		}); err != nil {
			return err
		}
		ts = ts.Add(time.Hour * -6)
	}
	return nil
}
