package worker

import (
	"os/exec"

	"github.com/takutakahashi/notion-tpl/pkg/notion"
)

type Worker struct {
	Client     notion.Client
	exportPath string
	tmplPath   string
	Cmd        string
}

func New(token, tbid, exportPath, tmplPath, cmd string) Worker {
	cli := notion.NewClient(token, tbid, exportPath)
	return Worker{
		Client:     cli,
		exportPath: exportPath,
		tmplPath:   tmplPath,
		Cmd:        cmd,
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
	exec.Command(w.Cmd).Run()
	return nil
}
