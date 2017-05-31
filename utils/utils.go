package utils

import (
	"html/template"
	"io"
)

func Tmpl(text string, data interface{}, wr io.Writer) error {

	t := template.New("Usage")
	template.Must(t.Parse(text))

	return t.Execute(wr, data)
}