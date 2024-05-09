package xtm

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

type TemplateManager struct {
	funcs             template.FuncMap
	baseTemplatesPath []string
	htmlTemplatesPath []string
	debug             bool
	errorTemplateName string

	baseTemplates *template.Template
	htmlTemplates map[string]*template.Template
}

func New(bases []string, htmls []string, funcs template.FuncMap) *TemplateManager {
	tm := &TemplateManager{
		funcs:             funcs,
		baseTemplatesPath: bases,
		htmlTemplatesPath: htmls,
		debug:             false,
		errorTemplateName: "error.tmpl",
	}
	return tm.Reload()
}

func NewGlob(basePattern string, htmlPattern string, funcs template.FuncMap) *TemplateManager {
	bases, err := filepath.Glob(basePattern)
	if err != nil {
		panic(err)
	}

	htmlx, err := filepath.Glob(htmlPattern)
	if err != nil {
		panic(err)
	}

	// 排除 base templates
	exclude := make(map[string]bool)
	for _, path := range bases {
		exclude[base(path)] = true
	}
	htmls := make([]string, 0)
	for _, path := range htmlx {
		if !exclude[base(path)] {
			htmls = append(htmls, path)
		}
	}

	return New(bases, htmls, funcs)
}

// SetDebug
// debug = true 则每次执行模板都会重新加载模板文件，修改模板文件后，不需要重启程序就会生效，调试用
// debug = false 只会在程序启动时加载一次模板文件，修改模板文件后，需要重启程序才会生效
func (tm *TemplateManager) SetDebug(debug bool) *TemplateManager {
	tm.debug = debug
	return tm
}

// SetErrorTemplateName
// 执行模板失败时会返回此模板
func (tm *TemplateManager) SetErrorTemplateName(name string) *TemplateManager {
	tm.errorTemplateName = name
	return tm
}

// Reload 加载所有模板
// htmlTemplates 存储子模板，baseTemplates 存储基础模板
// 子模板都相互独立，但每个子模板中都包含基础模板
// 这样做的目的是一个 Template 对象不能包含相同的模板名称，不利于子模板继承基础模板的设计模式
// 子模板名称为模板文件名称
func (tm *TemplateManager) Reload() *TemplateManager {
	baseTemplates := template.New("base.tmpl")
	baseTemplates.Funcs(tm.funcs)
	_, _ = baseTemplates.ParseFiles(tm.baseTemplatesPath...)

	htmlTemplates := make(map[string]*template.Template)
	for _, html := range tm.htmlTemplatesPath {
		// 以文件名为 key，key 不包含路径部分
		htmlTemplates[base(html)] = _must(_must(baseTemplates.Clone()).ParseFiles(html))
	}

	tm.baseTemplates = baseTemplates
	tm.htmlTemplates = htmlTemplates

	return tm
}

// Html 执行模板并将执行结果返回给客户端
// name 模板文件名称
// data 模板数据源
func (tm *TemplateManager) Html(wr http.ResponseWriter, name string, data any) error {
	if tm.debug {
		tm.Reload()
	}

	if pt, ok := tm.htmlTemplates[name]; !ok {
		return fmt.Errorf("not found template[%s]", name)
	} else {
		return pt.ExecuteTemplate(wr, name, data)
	}
}

// HtmlGin 适配 gin
func (tm *TemplateManager) HtmlGin(ctx *gin.Context, code int, name string, data any) {
	ctx.Status(code)
	err := tm.Html(ctx.Writer, name, data)
	if err == nil {
		return
	}

	errorCode := http.StatusInternalServerError

	ctx.Status(errorCode)
	err = tm.Html(ctx.Writer, tm.errorTemplateName, err.Error())
	if err != nil {
		ctx.String(errorCode, "%v", err)
	}
}

// PrintDefinedTemplates 输出所有模板名称，调试用
func (tm *TemplateManager) PrintDefinedTemplates() {
	names := tm.DefinedTemplates()
	for _, name := range names {
		fmt.Println(name)
	}
}

// DefinedTemplates 获取所有模板名称，调试用
func (tm *TemplateManager) DefinedTemplates() []string {
	names := make([]string, 0)

	names = append(names, tm.promoteDefinedTemplates(
		tm.baseTemplates.Name(),
		tm.baseTemplates.DefinedTemplates()),
	)

	for s, t := range tm.htmlTemplates {
		names = append(names, tm.promoteDefinedTemplates(s, t.DefinedTemplates()))
	}

	return names
}

// 优化模版名称显示
func (tm *TemplateManager) promoteDefinedTemplates(name string, defined string) string {
	defined = strings.ReplaceAll(defined, "\"", "")
	defined = strings.ReplaceAll(defined, "; defined templates are: ", "")

	segments1 := strings.Split(defined, ",")
	segments2 := make([]string, 0, len(segments1))
	for _, segment := range segments1 {
		segments2 = append(segments2, strings.TrimSpace(segment))
	}

	return name + " ==> " + strings.Join(segments2, ", ")
}
