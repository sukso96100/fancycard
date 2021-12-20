package tmpl_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
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
				fmt.Print("File reading error", err)
				return
			}
			fmt.Fprintln(w, string(data))
		}))
		defer ts.Close()
		templatePath, templateData, err := tmpl.ExtractMetaTagsFromURL(ts.URL)
		assert.Nil(t, err)
		assert.Equal(t, data["templatePath"].(string), templatePath)
		assert.Equal(t, reflect.DeepEqual(templateData, data["templateData"]), true)
	}
}

func Test_LoadTemplate_FromFile(t *testing.T) {
	templateFiles := []string{"simple.html", "simple2.html"}
	for _, templateFile := range templateFiles {
		templatePath := "templates/" + templateFile
		templateTextA, err := ioutil.ReadFile(templatePath)
		assert.Nil(t, err)
		templateTextB, err := tmpl.LoadTemplate(templateFile)
		assert.Nil(t, err)
		assert.Equal(t, string(templateTextA), templateTextB)
	}
}

func Test_LoadTemplate_FromURL(t *testing.T) {
	templateFiles := []string{"simple.html", "simple2.html"}
	for _, templateFile := range templateFiles {
		templatePath := "templates/" + templateFile
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			data, err := ioutil.ReadFile(templatePath)
			if err != nil {
				fmt.Println("File reading error", err)
				return
			}
			fmt.Fprint(w, string(data))
		}))
		defer ts.Close()
		templateTextA, err := ioutil.ReadFile(templatePath)
		assert.Nil(t, err)
		templateTextB, err := tmpl.LoadTemplate(ts.URL)
		assert.Nil(t, err)
		assert.Equal(t, string(templateTextA), templateTextB)
	}
}
