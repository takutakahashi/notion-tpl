package notion

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/kjk/notionapi"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/takutakahashi/notion-tpl/pkg/body"
)

type Client struct {
	c          *notionapi.Client
	TableID    string
	permMap    map[*notionapi.TableRow]time.Time
	exportPath string
	imagePath  string
}

func NewClient(token, tbid, exportPath, imagePath string) Client {
	client := &notionapi.Client{
		AuthToken: token,
	}

	return Client{
		c:          client,
		exportPath: exportPath,
		permMap:    map[*notionapi.TableRow]time.Time{},
		TableID:    tbid,
		imagePath:  imagePath,
	}
}

func (c Client) UpdatedPages() (map[*notionapi.Page]body.Body, error) {
	pages := map[*notionapi.Page]body.Body{}
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
		released := len(row.Columns[2]) != 0 && row.Columns[2][0].Text == "Yes"
		page, err := c.c.DownloadPage(row.Page.ID)
		if err != nil {
			return nil, err
		}
		pages[page] = body.New(page, row.Columns[1][0].Text, released)
	}
	return pages, nil
}

func (c Client) UploadImage(p *notionapi.Page) error {
	p.ForEachBlock(func(b *notionapi.Block) {
		if b.Type == notionapi.BlockImage {
			if err := c.uploadImageFromBlock(b); err != nil {
				fmt.Println(err)
			}
		}
	})
	return nil
}

func (c Client) uploadImageFromBlock(b *notionapi.Block) error {
	resp, err := c.c.DownloadFile(b.Source, b)
	if err != nil {
		return err
	}
	fmt.Println(b.Source)
	source := b.Source // also present in block.Format.DisplaySource
	// source looks like: "https://s3-us-west-2.amazonaws.com/secure.notion-static.com/e5470cfd-08f0-4fb8-8ec2-452ca1a3f05e/Schermafbeelding2018-06-19om09.52.45.png"
	var fileID string
	if len(b.FileIDs) > 0 {
		fileID = b.FileIDs[0]
	}
	parts := strings.Split(source, "/")
	fileName := parts[len(parts)-1]
	parts = strings.SplitN(fileName, ".", 2)
	ext := ""
	if len(parts) == 2 {
		fileName = parts[0]
		ext = "." + parts[1]
	}
	file := fmt.Sprintf("%s/%s-%s%s", c.imagePath, fileName, fileID, ext)
	if _, err := os.Stat(file); err != nil {
		return ioutil.WriteFile(fmt.Sprintf("%s/%s-%s%s", c.imagePath, fileName, fileID, ext), resp.Data, 0644)
	}
	return nil
}

func (c Client) LastUpdated() (time.Time, error) {
	logrus.Info(c.exportPath)
	files, err := ioutil.ReadDir(c.exportPath)
	ret := time.Unix(0, 0)
	if err != nil {
		return time.Time{}, errors.Wrapf(err, "failed to ReadDir")
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
