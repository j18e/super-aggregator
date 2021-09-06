package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/j18e/super-aggregator/models"
)

type Entry struct {
	Timestamp   string `json:"timestamp"`
	LogLine     string `json:"log_line"`
	Application string `json:"application"`
	Host        string `json:"host"`
	Environment string `json:"environment"`
	IP          string `json:"ip_address"`
}

func NewEntryController(es models.EntryService) EntryController {
	return EntryController{es}
}

type EntryController struct {
	es models.EntryService
}

func (ec *EntryController) EntriesHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		entries, err := ec.es.Entries()
		if err != nil {
			c.String(http.StatusInternalServerError, "something went wrong")
			return
		}
		c.HTML(http.StatusOK, "entries.html", entries)
	}
}

func (ec *EntryController) CreateHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var e Entry
		if err := c.ShouldBindJSON(e); err != nil {
			c.String(400, "%s", err)
			return
		}
		ts, err := time.Parse(time.RFC3339, e.Timestamp)
		if err != nil {
			c.String(400, "timestamp field must be formatted according to RFC3339")
			return
		}
		err = ec.es.Create(models.Entry{
			Timestamp:   ts,
			LogLine:     e.LogLine,
			Application: e.Application,
			Host:        e.Host,
			Environment: e.Environment,
			IP:          e.IP,
		})
		if err != nil {
			c.String(http.StatusInternalServerError, "internal server error")
			return
		}
	}
}
