package notion

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/jomei/notionapi"
	"github.com/nisanthchunduru/notion2markdown"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/takutakahashi/notion-tpl/pkg/body"
)

type Client struct {
	c          *notionapi.Client
	md         notion2markdown.Notion2Markdown
	TableID    string
	exportPath string
	imagePath  string
}

func NewClient(token, tbid, exportPath, imagePath string) Client {
	client := notionapi.NewClient(notionapi.Token(token))
	notion2Markdown := notion2markdown.Notion2Markdown{
		NotionToken: token,
	}

	return Client{
		c:          client,
		md:         notion2Markdown,
		exportPath: exportPath,
		TableID:    tbid,
		imagePath:  imagePath,
	}
}

func (c Client) UpdatedPages() (map[*notionapi.Page]body.Body, error) {
	ctx := context.Background()
	pages := map[*notionapi.Page]body.Body{}
	rows, err := c.c.Database.Query(ctx, notionapi.DatabaseID(c.TableID), nil)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to query database")
	}
	for _, row := range rows.Results {
		page, err := c.c.Page.Get(ctx, notionapi.PageID(row.ID))
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get page")
		}
		content, err := c.md.PageToMarkdown(page.ID.String())
		if err != nil {
			return nil, errors.Wrapf(err, "failed to convert page to markdown")
		}
		pages[page] = body.New(ctx, page, content)
	}
	return pages, nil
}

func (c Client) UploadImage(p *notionapi.Page) error {
	ctx := context.Background()
	blocks, err := c.c.Block.GetChildren(ctx, notionapi.BlockID(p.ID), nil)
	if err != nil {
		return errors.Wrapf(err, "failed to get blocks")
	}
	for _, block := range blocks.Results {
		if block.GetType() == notionapi.BlockTypeImage {
			if err := c.uploadImageFromBlock(block); err != nil {
				logrus.Error(err)
			}
		}
	}
	return nil

}

func (c Client) uploadImageFromBlock(b notionapi.Block) error {
	// download image
	srcURL := b.(notionapi.ImageBlock).Image.GetURL()
	res, err := http.Get(srcURL)
	if err != nil {
		return errors.Wrapf(err, "failed to get image")
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return errors.Wrapf(err, "failed to read image")
	}

	// source looks like: "https://s3-us-west-2.amazonaws.com/secure.notion-static.com/e5470cfd-08f0-4fb8-8ec2-452ca1a3f05e/Schermafbeelding2018-06-19om09.52.45.png"
	var fileID string
	parts := strings.Split(srcURL, "/")
	fileName := parts[len(parts)-1]
	parts = strings.SplitN(fileName, ".", 2)
	ext := ""
	if len(parts) == 2 {
		fileName = parts[0]
		ext = "." + parts[1]
	}
	file := fmt.Sprintf("%s/%s-%s%s", c.imagePath, fileName, ext)
	if _, err := os.Stat(file); err != nil {
		return os.WriteFile(fmt.Sprintf("%s/%s-%s%s", c.imagePath, fileName, fileID, ext), data, 0644)
	}
	return nil
}

func (c Client) LastUpdated() (time.Time, error) {
	logrus.Info(c.exportPath)
	files, err := os.ReadDir(c.exportPath)
	ret := time.Unix(0, 0)
	if err != nil {
		return time.Time{}, errors.Wrapf(err, "failed to ReadDir")
	}
	for _, file := range files {
		i, err := file.Info()
		if err != nil {
			return time.Time{}, errors.Wrapf(err, "failed to get file info")
		}
		mt := i.ModTime()
		if mt.After(ret) {
			ret = mt
		}
	}
	// Notion's time resolution is minute. So it needs to round.
	return ret.Add(-1 * time.Minute), nil
}
