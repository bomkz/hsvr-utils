package apifrontend

import "github.com/lxzan/gws"

var reconnecting bool

var success = make(chan bool)

var WsStreamClosed = make(chan bool)

var localSocket *gws.Conn
