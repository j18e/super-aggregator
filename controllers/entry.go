package controllers

import (
	"net/http"
	"net/url"
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
		var params struct {
			Application string `form:"application"`
			Host        string `form:"host"`
			Environment string `form:"environment"`
		}

		if err := c.ShouldBindQuery(&params); err != nil {
			c.String(http.StatusBadRequest, "%s", err)
			return
		}
		var data struct {
			Entries      []models.Entry
			Applications []namePath
			Hosts        []namePath
			Environments []namePath
		}
		eq := models.EntriesQuery{
			Application: params.Application,
			Host:        params.Host,
			Environment: params.Environment,
		}
		entries, err := ec.es.Entries(eq)
		if err != nil {
			c.String(http.StatusInternalServerError, "something went wrong")
			return
		}

		data.Entries = entries

		reqPath := c.Request.URL.Path
		apps, _ := ec.es.Applications()
		data.Applications = assembleDropdownData(c.Request.URL.Query(), apps, "application", reqPath)
		hosts, _ := ec.es.Hosts()
		data.Hosts = assembleDropdownData(c.Request.URL.Query(), hosts, "host", reqPath)
		envs, _ := ec.es.Environments()
		data.Environments = assembleDropdownData(c.Request.URL.Query(), envs, "environment", reqPath)

		c.HTML(http.StatusOK, "entries.html", data)
	}
}

type namePath struct {
	Name, Path string
}

func assembleDropdownData(form url.Values, data []string, field, path string) []namePath {
	form.Del(field)
	u := &url.URL{
		Path:     path,
		RawQuery: form.Encode(),
	}
	res := []namePath{{"All", u.String()}}
	for _, val := range data {
		form.Set(field, val)
		u := &url.URL{
			Path:     path,
			RawQuery: form.Encode(),
		}
		res = append(res, namePath{
			Name: val,
			Path: u.String(),
		})
	}
	return res
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
