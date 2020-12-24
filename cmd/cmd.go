package main

import (
	"log"
	"os"
	"time"

	"github.com/kjk/notionapi"
	"github.com/takutakahashi/notion-tpl/pkg/worker"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Name = "notion-tpl"
	app.Usage = "render template from notion page"
	app.Version = "0.0.1"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "table-id",
			Usage: "table view id. ex: a6b2dab9302744a2bcc4e00c3b512ae6",
		},
	}
	app.Action = action
	app.Run(os.Args)
}

type Body struct {
	Content   []byte
	Title     string
	Tags      []string
	UpdatedAt time.Time
}

func action(c *cli.Context) error {
	client := &notionapi.Client{
		AuthToken: os.Getenv("NOTION_TOKEN"),
	}
	tableID := c.String("table-id") // p, err := client.DownloadPage(row.Page.ID)
	// if err != nil {
	// 	return err
	// }
	// fmt.Printf("%s", tomarkdown.ToMarkdown(p))

	page, err := client.DownloadPage(tableID)
	if err != nil {
		log.Fatalf("DownloadPage() failed with %s\n", err)
	}
	tb := page.TableViews[0]
	w := worker.New(tb)
	w.Start()

	return nil
}
