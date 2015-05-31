package parse

import (
	"bytes"
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
			app.Application{},
		},
		{
			"pagination.sb",
			app.Application{},
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
			t.Error("Application not the same")
		}
	}
}
