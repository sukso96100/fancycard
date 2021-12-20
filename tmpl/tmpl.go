package tmpl

import (
	"bytes"
	"embed"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"strings"
	"text/template"

	"github.com/PuerkitoBio/goquery"
)

//go:embed templates/*.html
var builtInTemplates embed.FS

const metaNamePrefix = "fancycard:"

// Extracts template path and map of data from meta tags from given web url.
func ExtractMetaTagsFromURL(url string) (string, map[string][]string, error) {
	// Request the HTML page.
	res, err := http.Get(url)
	if err != nil {
		return "", nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return "", nil, err
	}

	templatePath := ""
	templateData := make(map[string][]string)

	// Find meta tags with metaNamePrefix and build map of data.
	doc.Find("meta").Each(func(i int, s *goquery.Selection) {
		if name, _ := s.Attr("name"); strings.HasPrefix(name, metaNamePrefix) {
			key := strings.TrimPrefix(name, metaNamePrefix)
			value, _ := s.Attr("content")
			if key == "template" && templatePath == "" {
				templatePath = value
			}
			if key != "template" {
				templateData[key] = append(templateData[key], value)
			}
		}
	})

	return templatePath, templateData, nil
}

// Loads template text from built in templates file or from URL.
func LoadTemplate(templatePath string) (string, error) {
	if strings.HasPrefix(templatePath, "http://") || strings.HasPrefix(templatePath, "https://") {
		resp, err := http.Get(templatePath)
		if err != nil {
			return "", err
		}

		defer resp.Body.Close()

		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		return string(data), nil
	}
	tmplText, err := builtInTemplates.ReadFile(path.Join("templates", templatePath))
	if err != nil {
		return "", err
	}
	return string(tmplText), nil
}

// Compile the given template text with data into complete text output.
func CompileTemplate(templateText string, data map[string][]string) (string, error) {
	tmpl, err := template.New("Template").Parse(templateText)
	if err != nil {
		return "", err
	}

	var compiled bytes.Buffer
	err = tmpl.Execute(&compiled, data)
	if err != nil {
		return compiled.String(), err
	}
	return compiled.String(), nil
}
