package body

import (
	"fmt"
	"time"

	"github.com/kjk/notionapi"
	"github.com/kjk/notionapi/tomarkdown"
)

type Body struct {
	Content   []byte
	Title     string
	Tags      []string
	permURI   string
	UpdatedAt time.Time
	CreatedAt time.Time
}

func New(page *notionapi.Page, permURI string) Body {
	fmt.Println(page.Root().Title)
	return Body{
		Content:   tomarkdown.ToMarkdown(page),
		Title:     page.Root().Title,
		permURI:   permURI,
		CreatedAt: page.Root().CreatedOn(),
		UpdatedAt: page.Root().LastEditedOn(),
	}
}
