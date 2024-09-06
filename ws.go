package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"log"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lxzan/gws"
)

func handleWS(message bytes.Buffer, datatype string) {
	DataTypeHandler(message, datatype)

}

func connectWS() {

	socket, _, err := gws.NewClient(new(WebSocket), &gws.ClientOption{
		Addr: "wss://hs.vtolvr.live",
	})

	if err != nil {
		log.Print(err)
		return
	}

	go socket.ReadLoop()

	for {
		var text = ""
		_, err := fmt.Scanf("%s", &text)
		if err != nil {
			return
		}
		if strings.TrimSpace(text) == "" {
			continue
		}
		err = socket.WriteString(text)
		if err != nil {
			return
		}
	}

}

func retryWS() {
	var recon = 0

	ReconnectTimer := time.NewTicker(30 * time.Second)
	reconnecting = true

	for {

		select {
		case <-ReconnectTimer.C:
			recon += 1

			go connectWS()

			log.Println("\nReconnection attempt " + strconv.Itoa(recon))
		case <-success:
			log.Println("\nReconnection attempt succeeded: Attempt #" + strconv.Itoa(recon))
			ReconnectTimer.Stop()
			reconnecting = false
			return
		}

	}

}

func (c *WebSocket) OnClose(_ *gws.Conn, err error) {
	log.Printf("onerror: err=%s\n", err.Error())
	if !reconnecting {
		wsStreamClosed <- true
		Socket = nil
		go retryWS()
	}
}

func (c *WebSocket) OnPong(_ *gws.Conn, _ []byte) {
}

func (c *WebSocket) OnOpen(socket *gws.Conn) {

	if reconnecting {
		success <- true
	}
	Socket = socket

	var Subscriptions SubscribeStruct

	Subscriptions.MessageType = "subscribe"
	Subscriptions.Data = append(Subscriptions.Data, "user_login", "user_logout", "death", "kill", "spawn", "online")

	subscriptionByte, err := json.Marshal(Subscriptions)
	if err != nil {
		log.Println(err)
		return
	}
	err = socket.WriteString(string(subscriptionByte))
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(string(subscriptionByte))
	log.Println("Client connection is open.")
}

func (c *WebSocket) OnPing(socket *gws.Conn, _ []byte) {
	var Pong PongStruct
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

	var DataType PayloadTypeStruct

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
