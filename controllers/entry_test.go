package controllers

import (
	"net/url"
	"reflect"
	"testing"
)

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
