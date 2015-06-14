package parse

import (
	"bytes"
	"fmt"
	"github.com/kr/pretty"
	"github.com/romanoff/servicebuilder/app"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"testing"
)

func TestParser(t *testing.T) {
	cases := []struct {
		Input       string
		Application app.Application
	}{
		{
			"simple.sb",
			app.Application{
				Models: []*app.Model{
					{
						Name:   "User",
						Fields: []*app.Field{{Name: "name", Type: app.STRING}},
					},
				},
			},
		},
		{
			"pagination.sb",
			app.Application{
				Models: []*app.Model{
					{
						Name:       "User",
						Fields:     []*app.Field{{Name: "name", Type: app.STRING}},
						Pagination: &app.Pagination{PerPage: 20, MaxPerPage: 100},
					},
				},
			},
		},
		{
			"actions.sb",
			app.Application{
				Models: []*app.Model{
					{
						Name:    "User",
						Fields:  []*app.Field{},
						Actions: &app.RestfulActions{Index: true, Show: true},
					},
				},
			},
		},
	}
	for _, tc := range cases {
		content, err := ioutil.ReadFile(filepath.Join("test-fixtures", tc.Input))
		if err != nil {
			t.Errorf("Expected to read content from %v fixture, but got %v", tc.Input, err)
		}
		parser := NewParser(tc.Input, bytes.NewBuffer([]byte(content)))
		application, err := parser.Parse()
		if err != nil {
			t.Errorf("Expected to not get error while parsing %v fixture, but got %v", tc.Input, err)
		}
		if !reflect.DeepEqual(application, tc.Application) {
			fmt.Printf("%# v\n", pretty.Formatter(application))
			fmt.Printf("%# v\n", pretty.Formatter(tc.Application))
			t.Error("Application not the same")
		}
	}
}
