package body

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"github.com/Masterminds/sprig"
	"github.com/kjk/notionapi"
	"github.com/kjk/notionapi/tomarkdown"
)

type Body struct {
	Content   string
	Title     string
	Tags      []string
	permURI   string
	Released  bool
	UpdatedAt time.Time
	CreatedAt time.Time
}

func New(page *notionapi.Page, permURI string, released bool) Body {
	return Body{
		Content:   fmt.Sprintf("%s", tomarkdown.ToMarkdown(page)),
		Title:     page.Root().Title,
		permURI:   permURI,
		CreatedAt: page.Root().CreatedOn(),
		UpdatedAt: page.Root().LastEditedOn(),
		Released:  released,
	}
}

func (p Body) Export(tmplPath, exportPath string) error {
	tmpl, err := template.New(filepath.Base(tmplPath)).Funcs(sprig.TxtFuncMap()).ParseFiles(tmplPath)
	if err != nil {
		return err
	}
	result, err := os.Create(fmt.Sprintf("%s/%s.md", exportPath, p.permURI))
	if err != nil {
		return err
	}
	defer result.Close()
	return tmpl.Execute(result, p)
}
