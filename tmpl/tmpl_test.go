package tmpl_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/sukso96100/fancycard/tmpl"
)

func Test_ExtractMetaTagsFromURL(t *testing.T) {
	testData := make([]map[string]interface{}, 0)
	testData = append(testData,
		map[string]interface{}{
			"testData":     "../testdata/page1.html",
			"templatePath": "simple.html",
			"templateData": map[string][]string{
				"Title": {"Test Target Page 1"},
			},
		})
	testData = append(testData,
		map[string]interface{}{
			"testData":     "../testdata/page2.html",
			"templatePath": "simple2.html",
			"templateData": map[string][]string{
				"Title":       {"Test Target Page 2"},
				"Description": {"This page uses simple2.html"},
			},
		})
	for _, data := range testData {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			data, err := ioutil.ReadFile(data["testData"].(string))
			if err != nil {
				fmt.Println("File reading error", err)
				return
			}
			fmt.Fprintln(w, string(data))
		}))
		defer ts.Close()
		templatePath, templateData, err := tmpl.ExtractMetaTagsFromURL(ts.URL)
		if err != nil {
			t.Error(err)
		}
		if templatePath != data["templatePath"].(string) {
			t.Errorf("templatePath: expected %s, got %s", data["templatePath"].(string), templatePath)
		}
		if !reflect.DeepEqual(templateData, data["templateData"]) {
			t.Errorf("templateData: expected %v, got %v", data["templateData"], templateData)
		}
	}

}
