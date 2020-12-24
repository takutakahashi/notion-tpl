package main

import (
	"os"

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

func action(c *cli.Context) error {
	tableID := c.String("table-id")
	token := os.Getenv("NOTION_TOKEN")

	w := worker.New(token, tableID)
	w.Start()

	return nil
}
