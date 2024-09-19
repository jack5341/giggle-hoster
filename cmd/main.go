package main

import (
	"github.com/jack5341/giggle-hoster/api/handler"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func main() {
	app := pocketbase.New()

	RegisterHooks(app)

	if err := app.Start(); err != nil {
		panic(err)
	}
}

func RegisterHooks(app *pocketbase.PocketBase) {
	app.OnRecordBeforeCreateRequest().Add(func(e *core.RecordCreateEvent) error {
		if e.Record.Collection().Name == "nodes" {
			if err := handler.CreateNode(e); err != nil {
				return err
			}

			if err := app.Dao().SaveRecord(e.Record); err != nil {
				return err
			}
		}
		return nil
	})
}
