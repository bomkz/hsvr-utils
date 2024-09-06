package apifrontend

import (
	"fmt"
	"log"

	"github.com/lxzan/gws"
)

func main() {
	socket, _, err := gws.NewClient(new(WebSocket), &gws.ClientOption{
		Addr: "wss://hs.vtolvr.live",
		PermessageDeflate: gws.PermessageDeflate{
			Enabled:               true,
			ServerContextTakeover: true,
			ClientContextTakeover: true,
		},
	})
	if err != nil {
		log.Printf(err.Error())
		return
	}
	go socket.ReadLoop()
}

type WebSocket struct {
}

func (c *WebSocket) OnClose(socket *gws.Conn, err error) {
	fmt.Printf("onerror: err=%s\n", err.Error())
}

func (c *WebSocket) OnPong(socket *gws.Conn, payload []byte) {
}

func (c *WebSocket) OnOpen(socket *gws.Conn) {
	_ = socket.WriteString("hello, there is client")
}

func (c *WebSocket) OnPing(socket *gws.Conn, payload []byte) {
	_ = socket.WritePong(payload)
}

func (c *WebSocket) OnMessage(socket *gws.Conn, message *gws.Message) {
	defer message.Close()
	fmt.Printf("recv: %s\n", message.Data.String())
}
