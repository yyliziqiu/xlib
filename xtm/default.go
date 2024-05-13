package xtm

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

var Default *TemplateManager

func InitDefault(bases []string, htmls []string, funcs template.FuncMap) {
	Default = New(bases, htmls, funcs)
}

func InitDefaultGlob(basePattern string, htmlPattern string, funcs template.FuncMap) {
	Default = NewGlob(basePattern, htmlPattern, funcs)
}

func SetDebug(debug bool) *TemplateManager {
	return Default.SetDebug(debug)
}

func SetErrorTemplateName(name string) *TemplateManager {
	return Default.SetErrorTemplateName(name)
}

func Reload() *TemplateManager {
	return Default.Reload()
}

func Html(wr http.ResponseWriter, name string, data any) error {
	return Default.Html(wr, name, data)
}

func HtmlGin(ctx *gin.Context, code int, name string, data any) {
	Default.HtmlGin(ctx, code, name, data)
}

func PrintDefinedTemplates() {
	Default.PrintDefinedTemplates()
}

func DefinedTemplates() []string {
	return Default.DefinedTemplates()
}
