package templateutil

import (
	"bytes"
	"text/template"
)

func Generate(tplStr string, data any) (string, error) {
	tmpl, err := template.New("tpl").Parse(tplStr)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
