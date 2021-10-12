package tmpl

import (
	"bytes"
	"embed"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
	"text/template"
)

//go:embed templates/*.html
var builtInTemplates embed.FS

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
