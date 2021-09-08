package controllers

import (
	"net/url"
	"reflect"
	"testing"
	"time"
)

func Test_parseDuration(t *testing.T) {
	for _, tt := range []struct {
		cnt  int
		unit string
		exp  time.Duration
	}{
		{1, "h", time.Hour},
		{1, "d", time.Hour * 24},
		{2, "w", time.Hour * 24 * 14},
	} {
		got, err := parseDuration(tt.cnt, tt.unit)
		if err != nil {
			t.Errorf("got error: %s", err)
		}
		if got != tt.exp {
			t.Errorf("expected %v, got %v", tt.exp, got)
		}
	}
}

func Test_assembleDropdownData(t *testing.T) {
	exp := []namePath{
		{"All", "/somepath?var1=someval"},
		{"aaa", "/somepath?var1=someval&var2=aaa"},
		{"bbb", "/somepath?var1=someval&var2=bbb"},
		{"ccc", "/somepath?var1=someval&var2=ccc"},
	}

	form := url.Values{
		"var1": []string{"someval"},
	}
	data := []string{
		"aaa",
		"bbb",
		"ccc",
	}
	got := assembleDropdownData(form, data, "var2", "/somepath")

	if !reflect.DeepEqual(got, exp) {
		t.Errorf("got %v, expected %v", got, exp)
	}
}
