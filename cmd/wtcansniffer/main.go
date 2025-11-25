package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"net/http"

	"github.com/gorilla/websocket"
	"go.einride.tech/can/pkg/socketcan"
)

type CanMessage struct {
	Id 		uint32	`json:"id"`
	Length 	uint8	`json:"length"`
	Data 	[]byte	`json:"data"`
}

type WtWebSocket struct {
	conn		*websocket.Conn
	addr		string
	upgrader 	websocket.Upgrader
	// mutex 	sync.Mutex
}

// TODO: Reconfigure these later to the settings file
const (
	CAN_CHANNEL			= "vcan0"
	WEBSOCKET_PORT		= "localhost:8080"
)

var (
	canMsgList = []CanMessage{}
	addr       = flag.String("addr", WEBSOCKET_PORT, "http service address")
	// upgrader   = websocket.Upgrader{
	// 	ReadBufferSize:  1024,
	// 	WriteBufferSize: 1024,
	// }
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
)

func main () {
	connCAN, err := handleCanBusConnect()
	if err != nil {
		fmt.Println("[ERROR] Unable to connect with CAN Bus: ", err)
	}
	defer connCAN.Close()

	// go readCanBus(connCAN)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handleWebSocket(w, r, connCAN)
	})

	fmt.Println("Listening for WebSocket request on provided port: " + WEBSOCKET_PORT)
	// NOTE:
	// - Starts http server listening on `addr`
	// - `ListenAndServer` will continue listening in the background after WebSocket disconnections, to try reconnect itself when it receives another `/ws` request from WebClient
	
	//err = http.ListenAndServe(*addr, nil)
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("[ERROR] ListenAndServer", err)
	}
}

func handleCanBusConnect() (net.Conn, error) {
	connCAN, err := socketcan.DialContext(context.Background(), "can", CAN_CHANNEL)
	if err != nil {
		return nil, fmt.Errorf("CAN Bus connection error: %w", err)
	}
	
	fmt.Println("[INFO] Successfully connected to CAN Bus")
	return connCAN, nil
}

func handleWebSocket(w http.ResponseWriter, r *http.Request, connCAN net.Conn) {
	connWS := WtWebSocket{}
	var err error

	connWS.conn, err = upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("[ERROR] Setting up WebSocket")
		return
	}
	defer connWS.conn.Close()
	fmt.Println("[INFO] Connected to WebSocket")


	// Now read CAN Bus, and send to client
	recv := socketcan.NewReceiver(connCAN)

	for recv.Receive() {
		frame := recv.Frame()

		newFrame := CanMessage{
			Id: frame.ID,
			Length: frame.Length,
			Data: frame.Data[:],
		}

		canMsgList = append(canMsgList, newFrame)

		json, err := json.Marshal(newFrame)
		if err != nil {
			fmt.Println("[ERROR] Marshalling JSON: ", err)
		}

		err = connWS.conn.WriteMessage(websocket.TextMessage, json)
		if err != nil {
			fmt.Println("[ERROR] Sending data to client: ", err)
			return
		}
	}
}

func readCanBus(connCAN net.Conn) {
	recv := socketcan.NewReceiver(connCAN)

	for recv.Receive() {
		fmt.Println("reeeeee")
		frame := recv.Frame()

		newFrame := CanMessage{
			Id: frame.ID,
			Length: frame.Length,
			Data: frame.Data[:],
		}

		canMsgList = append(canMsgList, newFrame)
	}
}
