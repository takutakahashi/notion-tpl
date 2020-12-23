package main

import (
	"fmt"
	"log"
	"os"

	"github.com/kjk/notionapi"
	"github.com/kjk/notionapi/tomarkdown"
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

func action(c *cli.Context) error {
	client := &notionapi.Client{
		AuthToken: os.Getenv("NOTION_TOKEN"),
	}
	tableID := c.String("table-id")
	page, err := client.DownloadPage(tableID)
	if err != nil {
		log.Fatalf("DownloadPage() failed with %s\n", err)
	}
	fmt.Println(page.TableViews[0].Rows[0].Page.GetTitle()[0].Text)
	for _, id := range page.Root().ContentIDs {
		fmt.Println(id)
	}
	fmt.Printf("%s", tomarkdown.ToMarkdown(page))
	return nil
}
