package gocoder

import (
	"bytes"
	"os"
	"text/template"
)

func TemplateFromFile(tmplPath string, env interface{}, fn template.FuncMap) (Codeable, error) {
	b, err := os.ReadFile(tmplPath)
	if err != nil {
		return nil, err
	}
	return Template(string(b), env, fn)
}

func Template(tmplContent string, env interface{}, fn template.FuncMap) (Codeable, error) {
	fOutput := bytes.NewBuffer(nil)
	t, err := template.New("").Funcs(fn).Parse(tmplContent)
	if err != nil {
		return nil, err
	}
	err = t.Execute(fOutput, env)
	if err != nil {
		return nil, err
	}
	return NewValue(fOutput.String(), nil), err
}
