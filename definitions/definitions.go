package definitions

import (
	_ "embed"

	"fyne.io/fyne/v2"
	"github.com/google/uuid"
)

//go:embed .\..\aircraft.ico
var Icon []byte

type PongStruct struct {
	PID         uuid.UUID `json:"pid"`
	MessageType string    `json:"type"`
}
type PayloadTypeStruct struct {
	MessageType string `json:"type"`
}

var FrontendWindow fyne.Window
