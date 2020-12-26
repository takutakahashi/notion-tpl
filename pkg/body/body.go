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
	}
}

func (p Body) ExportHugo() error {
	path := "./src/hugo.md.tpl"
	tmpl, err := template.New(filepath.Base(path)).Funcs(sprig.TxtFuncMap()).ParseFiles(path)
	if err != nil {
		return err
	}
	result, err := os.Create(p.GetPathToExport())
	if err != nil {
		return err
	}
	defer result.Close()
	return tmpl.Execute(result, p)
}

func (p Body) GetPathToExport() string {
	return fmt.Sprintf("./content/posts/%s.md", p.permURI)
}
