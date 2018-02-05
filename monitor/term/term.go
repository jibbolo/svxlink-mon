package term

import (
	"time"

	ui "github.com/gizak/termui"
)

const label = " SVXLINK-MON - Press Q to quit - https://goo.gl/igcxWT "

type RowGetter interface {
	GetRows() [][]string
}

func Loop(src RowGetter) error {
	if err := ui.Init(); err != nil {
		return err
	}
	defer ui.Close()

	go func() {
		for {
			Render(src.GetRows())
			time.Sleep(100 * time.Millisecond)
		}
	}()

	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
	})

	ui.Handle("/sys/kbd/Q", func(ui.Event) {
		ui.StopLoop()
	})

	ui.Loop()
	return nil
}

func Close() {
	ui.Close()
}

func Render(rows [][]string) {
	rowsWithHeader := append([][]string{
		[]string{"ID", "IP Addr", "Status", "Time"},
	}, rows...)

	table := ui.NewTable()
	table.Rows = rowsWithHeader
	table.FgColor = ui.ColorWhite
	table.BgColor = ui.ColorDefault
	table.Y = 0
	table.BorderLabel = label
	table.BorderLabelFg = ui.ColorRed
	table.X = 0
	table.Analysis()
	table.SetSize()

	ui.Render(table)
}
