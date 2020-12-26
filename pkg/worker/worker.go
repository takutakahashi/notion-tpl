package worker

import (
	"fmt"

	"github.com/takutakahashi/notion-tpl/pkg/notion"
)

type Worker struct {
	Client notion.Client
}

func New(token, tbid string) Worker {
	cli := notion.NewClient(token, tbid, ".")
	return Worker{
		Client: cli,
	}
}

func (w Worker) Start() error {
	pages, err := w.Client.UpdatedPages()
	if err != nil {
		return err
	}
	for _, p := range pages {
		fmt.Println(p.Released)
		err = p.ExportHugo()
		if err != nil {
			return err
		}
	}
	return nil
}
