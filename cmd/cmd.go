package main

import (
	"log"
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
		cli.StringFlag{
			Name:  "token",
			Usage: "notion_v2 token",
		},
		cli.StringFlag{
			Name:  "export-path",
			Usage: "export path",
		},
		cli.StringFlag{
			Name:  "image-path",
			Usage: "image path",
		},
		cli.StringFlag{
			Name:  "template",
			Usage: "template file path",
		},
		cli.StringFlag{
			Name:  "cmd",
			Usage: "cmd that will be executed after update",
		},
		cli.BoolFlag{
			Name:  "once",
			Usage: "execute once",
		},
	}
	app.Action = action
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func action(c *cli.Context) error {
	tableID := c.String("table-id")
	token := c.String("token")
	exportPath := c.String("export-path")
	tmplPath := c.String("template")
	imagePath := c.String("image-path")
	cmd := c.String("cmd")
	once := c.Bool("once")
	w := worker.New(token, tableID, exportPath, tmplPath, imagePath, cmd)
	return w.Start(once)
}
