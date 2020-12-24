package worker

import (
	"time"

	"github.com/takutakahashi/notion-tpl/pkg/store"

	"github.com/kjk/notionapi"
)

type Worker struct {
	Table   *notionapi.TableView
	permMap map[*notionapi.TableRow]time.Time
	store   store.Store
}

func New(tb *notionapi.TableView) Worker {
	store := store.New()
	return Worker{
		store:   store,
		permMap: map[*notionapi.TableRow]time.Time{},
		Table:   tb,
	}
}

func (w Worker) Start() error {
	lastUpdate := w.store.LastUpdated()
	for _, row := range w.Table.Rows {
		w.permMap[row] = row.Page.LastEditedOn()
	}
	for row, v := range w.permMap {
		if lastUpdate.Before(v) {
			_ = row
		}
	}
	return nil
}
