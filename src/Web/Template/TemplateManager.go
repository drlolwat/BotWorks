package Template

import (
	"html/template"
)

var (
	tpl *template.Template
)

func Init() {
	tpl = template.Must(template.ParseGlob("./templates/*.html"))
}

func GetTpl() *template.Template {
	return tpl
}
