package apifrontend

import (
	"bytes"
	"encoding/json"

	"log"

	"github.com/bomkz/hsvr-utils/definitions"
	"github.com/google/uuid"
	"github.com/lxzan/gws"
)

func handleWS(message bytes.Buffer, datatype string) {

}

type WebSocket struct{}

func connectWS() error {

	socket, _, err := gws.NewClient(new(WebSocket), &gws.ClientOption{
		Addr: "wss://hs.vtolvr.live",
	})

	if err != nil {
		log.Print(err)
		return err
	}

	go socket.ReadLoop()
	return nil

}

func (c *WebSocket) OnClose(_ *gws.Conn, err error) {
	log.Printf("onerror: err=%s\n", err.Error())

}

func (c *WebSocket) OnPong(_ *gws.Conn, _ []byte) {
}

func (c *WebSocket) OnOpen(socket *gws.Conn) {

	localSocket = socket

	log.Println("Client connection is open.")
}

func (c *WebSocket) OnPing(socket *gws.Conn, _ []byte) {
	var Pong definitions.PongStruct
	Pong.MessageType = "pong"
	Pong.PID = uuid.New()
	pongByte, err := json.Marshal(Pong)
	if err != nil {
		log.Println(err)
		return
	}
	err = socket.WriteString(string(pongByte))
	if err != nil {
		log.Println(err)
		return
	}
}

func (c *WebSocket) OnMessage(socket *gws.Conn, message *gws.Message) {
	defer func(message *gws.Message) {
		err := message.Close()
		if err != nil {
			log.Println(err)
			return
		}
	}(message)

	var DataType definitions.PayloadTypeStruct

	err := json.Unmarshal(message.Bytes(), &DataType)
	if err != nil {
		log.Println(err)
		return
	}

	if DataType.MessageType == "ping" {
		c.OnPing(socket, nil)
		return
	}
	handleWS(*message.Data, DataType.MessageType)
}
