package worker

import (
	"github.com/takutakahashi/notion-tpl/pkg/notion"
)

type Worker struct {
	Client     notion.Client
	exportPath string
	tmplPath   string
}

func New(token, tbid, exportPath, tmplPath string) Worker {
	cli := notion.NewClient(token, tbid, exportPath)
	return Worker{
		Client:     cli,
		exportPath: exportPath,
		tmplPath:   tmplPath,
	}
}

func (w Worker) Start() error {
	pages, err := w.Client.UpdatedPages()
	if err != nil {
		return err
	}
	for _, p := range pages {
		err = p.Export(w.tmplPath, w.exportPath)
		if err != nil {
			return err
		}
	}
	return nil
}
