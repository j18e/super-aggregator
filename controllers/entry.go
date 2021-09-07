package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/j18e/super-aggregator/models"
)

type Entry struct {
	Timestamp   string `json:"timestamp"   binding:"required"`
	LogLine     string `json:"log_line"    binding:"required"`
	Application string `json:"application" binding:"required,alphanum"`
	Host        string `json:"host"        binding:"required,alphanum"`
	Environment string `json:"environment" binding:"required,alphanum"`
}

func NewEntryController(es models.EntryService) EntryController {
	return EntryController{es}
}

type EntryController struct {
	es models.EntryService
}

func (ec *EntryController) EntriesHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var qs struct {
			Application string `form:"application"`
		}
		if err := c.ShouldBindQuery(&qs); err != nil {
			c.String(http.StatusBadRequest, "%s", err)
			return
		}
		var data struct {
			Entries      []models.Entry
			Applications []string
		}
		eq := models.EntriesQuery{Application: qs.Application}
		entries, err := ec.es.Entries(eq)
		if err != nil {
			c.String(http.StatusInternalServerError, "something went wrong")
			return
		}
		data.Entries = entries
		apps, _ := ec.es.Applications()
		data.Applications = apps
		c.HTML(http.StatusOK, "entries.html", data)
	}
}

func (ec *EntryController) CreateHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var e Entry
		if err := c.ShouldBindJSON(&e); err != nil {
			c.String(400, "%s", err)
			return
		}
		ts, err := time.Parse(time.RFC3339, e.Timestamp)
		if err != nil {
			c.String(400, "timestamp field must be formatted according to RFC3339")
			return
		}
		ip, _ := c.RemoteIP()
		err = ec.es.Create(models.Entry{
			Timestamp:   ts,
			LogLine:     e.LogLine,
			Application: e.Application,
			Host:        e.Host,
			Environment: e.Environment,
			IP:          ip.String(),
		})
		if err != nil {
			c.String(http.StatusInternalServerError, "internal server error")
			return
		}
	}
}
