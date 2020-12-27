package worker

import (
	"os/exec"
	"time"

	"github.com/sirupsen/logrus"
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
	for {
		if err := w.execute(); err != nil {
			logrus.Error(err)
		}
		time.Sleep(1 * time.Minute)
	}
}

func (w Worker) execute() error {
	pages, err := w.Client.UpdatedPages()
	if err != nil {
		return err
	}
	for _, p := range pages {

		err = p.Export(w.tmplPath, w.exportPath)
		if err != nil {
			return err
		}
		logrus.Info("Updated.", "post=", p.Title)
	}
	return exec.Command(w.Cmd).Run()
}
