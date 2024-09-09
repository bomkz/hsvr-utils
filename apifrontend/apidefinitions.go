package apifrontend

import (
	"bytes"

	"github.com/lxzan/gws"
)

type messageStruct struct {
	message bytes.Buffer
	PID     string
}

var WsStreamClosed = make(chan bool)

var localSocket *gws.Conn

var messages []messageStruct
