package notion

import (
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/kjk/notionapi"
	"github.com/takutakahashi/notion-tpl/pkg/body"
)

type Client struct {
	c          *notionapi.Client
	Table      *notionapi.TableView
	permMap    map[*notionapi.TableRow]time.Time
	exportPath string
}

func NewClient(token, tbid, exportPath string) Client {
	client := &notionapi.Client{
		AuthToken: token,
	}
	page, err := client.DownloadPage(tbid)
	if err != nil {
		log.Fatalf("DownloadPage() failed with %s\n", err)
	}
	tb := page.TableViews[0]
	return Client{
		c:          client,
		exportPath: exportPath,
		permMap:    map[*notionapi.TableRow]time.Time{},
		Table:      tb,
	}
}

func (c Client) UpdatedPages() ([]body.Body, error) {
	// TODO: use query
	lastUpdated, err := c.LastUpdated()
	if err != nil {
		return nil, err
	}
	fmt.Println(lastUpdated)
	pages := []body.Body{}
	for _, row := range c.Table.Rows {
		c.permMap[row] = row.Page.LastEditedOn()
	}
	for row, v := range c.permMap {
		fmt.Println(v)
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
