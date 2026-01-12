package appli

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

func OuvertureApp() {
	a := app.New()
	w := a.NewWindow("Groupie Tracker")

	artists := Api()

	var items []string
	for _, a := range artists {
		items = append(items, fmt.Sprintf("%s - %s", a.Name, a.FirstAlbum))
	}

	list := widget.NewList(
		func() int { return len(items) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			obj.(*widget.Label).SetText(items[id])
		},
	)

	w.SetContent(list)
	w.ShowAndRun()
}
