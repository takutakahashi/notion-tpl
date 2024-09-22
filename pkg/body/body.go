package body

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/jomei/notionapi"
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

func New(ctx context.Context, page *notionapi.Page, content string) Body {
	var title, permURI string
	var released bool
	for name, prop := range page.Properties {
		switch name {
		case "Name":
			title = prop.(*notionapi.TitleProperty).Title[0].PlainText
		case "Tags":
			continue
		case "release":
			released = prop.(*notionapi.CheckboxProperty).Checkbox
		case "Permanent URL":
			permURI = prop.(*notionapi.TextProperty).Text[0].PlainText
		}
	}
	return Body{
		Content:   strings.TrimLeft(content, fmt.Sprintf("# %s\n", title)),
		Title:     title,
		permURI:   permURI,
		CreatedAt: page.CreatedTime,
		UpdatedAt: page.LastEditedTime,
		Released:  released,
	}
}

func (p Body) Export(tmplPath, exportPath string) error {
	tmpl, err := template.New(filepath.Base(tmplPath)).ParseFiles(tmplPath)
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
