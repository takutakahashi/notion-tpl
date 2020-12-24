package worker

import (
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
	_ = pages
	return nil
}
