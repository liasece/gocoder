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

func TemplateRowFromFile(tmplPath string, env interface{}, fn template.FuncMap) (string, error) {
	b, err := os.ReadFile(tmplPath)
	if err != nil {
		return "", err
	}
	return TemplateRaw(string(b), env, fn)
}

func Template(tmplContent string, env interface{}, fn template.FuncMap) (Codeable, error) {
	raw, err := TemplateRaw(tmplContent, env, fn)
	return NewValue(raw, nil), err
}

func TemplateRaw(tmplContent string, env interface{}, fn template.FuncMap) (string, error) {
	fOutput := bytes.NewBuffer(nil)
	t, err := template.New("").Funcs(fn).Parse(tmplContent)
	if err != nil {
		return "", err
	}
	err = t.Execute(fOutput, env)
	if err != nil {
		return "", err
	}
	return fOutput.String(), err
}
