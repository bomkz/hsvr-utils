package apifrontend

import (
	"encoding/json"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/bomkz/vtolvr-utils/definitions"
	"github.com/google/uuid"
)

func BuildFrontend() {
	input := widget.NewEntry()
	input.SetPlaceHolder("Enter User ID")

	content := container.New(layout.NewVBoxLayout(), input, widget.NewButton("Query", func() { queryUser(input.Text) }))

	//content := container.NewVBox(input, widget.NewButton("Query", func() { log.Println("Query was ", input.Text) }))
	definitions.FrontendWindow.SetContent(content)
	s := fyne.Size{
		Height: 512,
		Width:  512,
	}

	definitions.FrontendWindow.SetCloseIntercept(PreventFrontendClose)

	definitions.FrontendWindow.Resize(s)
	definitions.FrontendWindow.Show()

}

func PreventFrontendClose() {
	definitions.FrontendWindow.Hide()
}

func queryUser(UID string) {

	newLookup := definitions.UserLookup{
		MessageType: "lookup",
		Data: definitions.UserLookupData{
			UID:      UID,
			PID:      uuid.NewString(),
			Category: "user",
		},
	}

	lookupBytes, err := json.Marshal(newLookup)
	if err != nil {
		log.Print(err)
		return
	}
	definitions.Socket.WriteString(string(lookupBytes))
}
