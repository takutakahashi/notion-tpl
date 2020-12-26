package notion

import (
	"log"
	"time"

	"github.com/kjk/notionapi"
	"github.com/takutakahashi/notion-tpl/pkg/body"
	"github.com/takutakahashi/notion-tpl/pkg/store"
)

type Client struct {
	c       *notionapi.Client
	Table   *notionapi.TableView
	permMap map[*notionapi.TableRow]time.Time
	store   store.Store
}

func NewClient(token, tbid, storePath string) Client {
	client := &notionapi.Client{
		AuthToken: token,
	}
	store := store.New(storePath)
	page, err := client.DownloadPage(tbid)
	if err != nil {
		log.Fatalf("DownloadPage() failed with %s\n", err)
	}
	tb := page.TableViews[0]
	return Client{
		c:       client,
		store:   store,
		permMap: map[*notionapi.TableRow]time.Time{},
		Table:   tb,
	}
}

func (c Client) UpdatedPages() ([]body.Body, error) {
	// TODO: use query
	lastUpdated := c.store.LastUpdated()
	defer c.store.RefreshUpdated()
	pages := []body.Body{}
	for _, row := range c.Table.Rows {
		c.permMap[row] = row.Page.LastEditedOn()
	}
	for row, v := range c.permMap {
		if v.After(lastUpdated) {
			released := len(row.Columns[2]) != 0 && row.Columns[2][0].Text == "Yes"
			page, err := c.c.DownloadPage(row.Page.ID)
			if err != nil {
				return nil, err
			}
			pages = append(pages, body.New(page, row.Columns[1][0].Text, released))
		}

	}
	return pages, nil
}
