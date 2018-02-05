package term

import (
	"time"

	ui "github.com/gizak/termui"
)

type rowgetter interface {
	GetRows() [][]string
}

func Render(rows [][]string) {
	rowsWithHeader := append([][]string{
		[]string{"ID", "IP Addr", "Status"},
	}, rows...)

	table := ui.NewTable()
	table.Rows = rowsWithHeader
	table.FgColor = ui.ColorWhite
	table.BgColor = ui.ColorDefault
	table.Y = 0
	table.X = 0
	table.Analysis()
	table.SetSize()

	ui.Render(table)
}

func Run(src rowgetter) error {
	err := ui.Init()
	if err != nil {
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
	ui.Handle("/sys/kbd/C-x", func(ui.Event) {
		ui.StopLoop()
	})
	ui.Loop()
	return nil
}
