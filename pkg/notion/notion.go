package notion

import (
	"errors"
	"io/ioutil"
	"log"
	"time"

	"github.com/kjk/notionapi"
	"github.com/takutakahashi/notion-tpl/pkg/body"
)

type Client struct {
	c          *notionapi.Client
	TableID    string
	permMap    map[*notionapi.TableRow]time.Time
	exportPath string
}

func NewClient(token, tbid, exportPath string) Client {
	client := &notionapi.Client{
		AuthToken: token,
	}

	return Client{
		c:          client,
		exportPath: exportPath,
		permMap:    map[*notionapi.TableRow]time.Time{},
		TableID:    tbid,
	}
}

func (c Client) UpdatedPages() ([]body.Body, error) {
	// TODO: use query
	lastUpdated, err := c.LastUpdated()
	if err != nil {
		return nil, err
	}
	pages := []body.Body{}
	page, err := c.c.DownloadPage(c.TableID)
	if err != nil {
		log.Fatalf("DownloadPage() failed with %s\n", err)
		return nil, err
	}
	if page == nil || len(page.TableViews) == 0 {
		return nil, errors.New("page was wrong data")
	}
	tb := page.TableViews[0]
	for _, row := range tb.Rows {
		if row.Page.LastEditedOn().After(lastUpdated) {
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

func (c Client) LastUpdated() (time.Time, error) {
	files, err := ioutil.ReadDir(c.exportPath)
	ret := time.Unix(0, 0)
	if err != nil {
		return time.Time{}, err
	}
	for _, file := range files {
		mt := file.ModTime()
		if mt.After(ret) {
			ret = mt
		}
	}
	// Notion's time resolution is minute. So it needs to round.
	return ret.Add(-1 * time.Minute), nil
}
