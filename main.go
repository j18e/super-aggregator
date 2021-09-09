package main

import (
	"errors"
	"flag"
	"fmt"
	"time"

	"github.com/j18e/super-aggregator/controllers"
	"github.com/j18e/super-aggregator/models"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
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
	dbDebug := flag.Bool("db.debug", false, "whether to log all database queries")
	pgHost := flag.String("pg.host", "localhost", "host running the postgres server")
	pgUser := flag.String("pg.user", "super-aggregator", "name of the user and database to connect to")
	pgPass := flag.String("pg.password", "", "password to authenticate to postgres with")
	pgPort := flag.Int("pg.port", 5432, "port number on which postgres is running")
	flag.Parse()

	gormCfg := &gorm.Config{}
	if *dbDebug {
		gormCfg.Logger = logger.Default.LogMode(logger.Info)
	}

	if *pgPass == "" {
		return errors.New("flag -pg.password required when using postgres driver")
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		*pgHost, *pgUser, *pgPass, *pgUser, *pgPort)
	gdb, err := gorm.Open(postgres.Open(dsn), gormCfg)
	if err != nil {
		return err
	}
	es := models.NewEntryService(gdb)

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
	if err := es.Create(create...); err != nil {
		return err
	}
	return nil
}
