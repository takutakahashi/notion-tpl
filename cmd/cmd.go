package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/kjk/notionapi"
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
	tableID := c.String("table-id")
	page, err := client.DownloadPage(tableID)
	if err != nil {
		log.Fatalf("DownloadPage() failed with %s\n", err)
	}
	tb := page.TableViews[0]
	permMap := map[string]time.Time{}
	for _, row := range tb.Rows {
		permMap[row.Columns[1][0].Text] = row.Page.LastEditedOn()
		fmt.Println(row.Page.LastEditedOn())
		// p, err := client.DownloadPage(row.Page.ID)
		// if err != nil {
		// 	return err
		// }
		// fmt.Printf("%s", tomarkdown.ToMarkdown(p))
	}
	fmt.Printf("%s", permMap)
	return nil
}
