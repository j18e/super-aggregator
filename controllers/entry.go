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

type queryParams struct {
	Application string `form:"application"`
	Host        string `form:"host"`
	Environment string `form:"environment"`
}

func (ec *EntryController) EntriesHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var params queryParams

		if err := c.ShouldBindQuery(&params); err != nil {
			c.String(http.StatusBadRequest, "%s", err)
			return
		}
		var data struct {
			Entries      []models.Entry
			Applications []namePath
			Hosts        []namePath
			Environments []namePath
			Current      queryParams
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
		if params.Application == "" {
			params.Application = "All"
		}
		if params.Host == "" {
			params.Host = "All"
		}
		if params.Environment == "" {
			params.Environment = "All"
		}
		data.Current = params

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
		var ex []Entry
		if err := c.ShouldBindJSON(&ex); err != nil {
			c.String(400, "%s", err)
			return
		}
		var create []models.Entry
		for _, e := range ex {
			ts, err := time.Parse(time.RFC3339, e.Timestamp)
			if err != nil {
				c.String(400, "timestamp field must be formatted according to RFC3339")
				return
			}
			ip, _ := c.RemoteIP()
			create = append(create, models.Entry{
				Timestamp:   ts,
				LogLine:     e.LogLine,
				Application: e.Application,
				Host:        e.Host,
				Environment: e.Environment,
				IP:          ip.String(),
			})
		}
		if err := ec.es.Create(create); err != nil {
			c.String(http.StatusInternalServerError, "internal server error")
			return
		}
	}
}
