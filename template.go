package main

import (
	"fmt"
	"html/template"
	"io"
)

type TemplatePointer struct {
	*template.Template
}

type ApiTemplate struct {
	Info      TemplatePointer
	Mindustry TemplatePointer
}

func (t TemplatePointer) WriteData(w io.Writer, data interface{}) {

	err := t.Execute(w, data)
	if err != nil {
		if _, e := w.Write([]byte(err.Error())); e != nil {
			fmt.Println(e)
		}
	}
}

func (t TemplatePointer) WriteError(w io.Writer, err error) {
	if _, e := w.Write([]byte(err.Error())); e != nil {
		fmt.Println(e)
	}
}
func initApiTemplate(viewDir string) (ApiTemplate, error) {
	tp, err := readApiTemplate(
		[]string{"info", "mdt"},
		viewDir+"/api")
	if err != nil {
		return ApiTemplate{}, err
	}

	return ApiTemplate{
		Info:      tp[0],
		Mindustry: tp[1],
	}, nil
}
func readApiTemplate(htmlFileName []string, viewDir string) ([]TemplatePointer, error) {
	var apiTemplate []TemplatePointer
	head := viewDir + "/head.gohtml"
	footer := viewDir + "/footer.gohtml"
	for _, name := range htmlFileName {
		tp, err := template.New(name+".gohtml").
			ParseFiles(viewDir+"/"+name+".gohtml", head, footer)
		if err != nil {
			return nil, err
		}
		apiTemplate = append(apiTemplate, TemplatePointer{tp})
	}
	return apiTemplate, nil
}
