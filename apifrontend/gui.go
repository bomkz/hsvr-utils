package apifrontend

import (
	"bytes"
	"encoding/json"
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/bomkz/hsvr-utils/definitions"
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

	newLookup := definitions.LookupStruct{
		MessageType: "lookup",
		Data: definitions.LookupDataStruct{
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

	err = connectWS()
	if err != nil {
		log.Fatal(err)
	}

	localSocket.WriteString(string(lookupBytes))

	timer := time.NewTicker(1 * time.Second)

	var result bytes.Buffer
	for range timer.C {
		found := false
		for _, y := range messages {
			if y.PID == newLookup.Data.PID {
				result = y.message
				found = true
				break
			}
		}
		if found {
			timer.Stop()
			break
		}
	}
	localSocket.NetConn().Close()
	log.Println(result.String())

}
