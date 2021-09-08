package main

import (
	"flag"
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
	testData := flag.Bool("test-data", false, "flush the database and generate new test data before starting the web server")
	flag.Parse()

	db, err := gorm.Open(sqlite.Open("tmp/data.sqlite"), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		return err
	}
	es := models.NewEntryService(db)

	if *testData {
		if err := es.DestructiveReset(); err != nil {
			return err
		}
		if err := createEntries(es); err != nil {
			return err
		}
	} else {
		if err := es.AutoMigrate(); err != nil {
			return err
		}
	}

	ec := controllers.NewEntryController(es)

	r := gin.Default()
	r.LoadHTMLGlob("./views/*.html")
	r.GET("/", ec.EntriesHandler())
	r.POST("/", ec.EntriesTimePickerHandler())
	r.POST("/api/entry", ec.CreateHandler())
	return r.Run(":9000")
}

func createEntries(es models.EntryService) error {
	ts := time.Now()
	var create []models.Entry
	for i := 0; i < 100; i++ {
		create = append(create, models.Entry{
			Timestamp:   ts,
			LogLine:     "something exciting happened",
			Application: "app1",
			Host:        "host1",
			Environment: "prod",
		})
		ts = ts.Add(time.Hour * -6)
		create = append(create, models.Entry{
			Timestamp:   ts,
			LogLine:     "something boring happened",
			Application: "app2",
			Host:        "host2",
			Environment: "prod",
		})
		ts = ts.Add(time.Hour * -6)
		create = append(create, models.Entry{
			Timestamp:   ts,
			LogLine:     "something whatever happened",
			Application: "app3",
			Host:        "host1",
			Environment: "test",
		})
		ts = ts.Add(time.Hour * -6)
		create = append(create, models.Entry{
			Timestamp:   ts,
			LogLine:     "nothing happened",
			Application: "app4",
			Host:        "host2",
			Environment: "test",
		})
		ts = ts.Add(time.Hour * -6)
	}
	if err := es.Create(create); err != nil {
		return err
	}
	return nil
}
