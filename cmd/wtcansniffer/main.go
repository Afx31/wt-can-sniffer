package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
	"go.einride.tech/can/pkg/socketcan"
)

// TODO: Trialing2
type FrameUpdate struct {
	Id     uint32  `json:"id"`
	Length uint8   `json:"length"`
	Data   [8]byte `json:"data"`
}

type CanData struct {
	Byte0 uint8 `json:"byte0"`
	Byte1 uint8 `json:"byte1"`
	Byte2 uint8 `json:"byte2"`
	Byte3 uint8 `json:"byte3"`
	Byte4 uint8 `json:"byte4"`
	Byte5 uint8 `json:"byte5"`
	Byte6 uint8 `json:"byte6"`
	Byte7 uint8 `json:"byte7"`
}

type CanFrame struct {
	Id     uint32    `json:"id"`
	Length uint8     `json:"length"`
	Data   []CanData `json:"data"`
}

type WtWebSocket struct {
	conn     *websocket.Conn
	addr     string
	upgrader websocket.Upgrader
	// mutex 	sync.Mutex
}

type WtWebSocketResponse struct {
	Type int    `json:"id"`
	Data string `json:"data"`
}

// TODO: Reconfigure these later to the settings file
const (
	CAN_CHANNEL    = "vcan0"
	WEBSOCKET_PORT = "localhost:8080"
)

var (
	addr = flag.String("addr", WEBSOCKET_PORT, "http service address")
	// upgrader   = websocket.Upgrader{
	// 	ReadBufferSize:  1024,
	// 	WriteBufferSize: 1024,
	// }
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	wg sync.WaitGroup

	// cfMap = make(map[uint32][]CanFrame)
	cfMap = make(map[uint32]CanFrame)
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


	// Listen for responses from client
	wg.Add(1)
	go func() {
		defer wg.Done()
		webSocketListener(connWS)
	}()


	// Now read CAN Bus, and send to client
	recv := socketcan.NewReceiver(connCAN)
	for recv.Receive() {
		frame := recv.Frame()

		incomingFrameID := frame.ID
		incomingCanData := CanData{
			Byte0: frame.Data[0],
			Byte1: frame.Data[1],
			Byte2: frame.Data[2],
			Byte3: frame.Data[3],
			Byte4: frame.Data[4],
			Byte5: frame.Data[5],
			Byte6: frame.Data[6],
			Byte7: frame.Data[7],
		}

		canFrame, exists := cfMap[incomingFrameID]
		if !exists {
			canFrame = CanFrame{
				Id:     incomingFrameID,
				Length: frame.Length,
				Data:   []CanData{},
			}
		}
		canFrame.Data = append(canFrame.Data, incomingCanData)
		cfMap[incomingFrameID] = canFrame

		frameUpdate := FrameUpdate{
			Id:     incomingFrameID,
			Length: frame.Length,
			Data:   frame.Data,
		}

		jsonData, err := json.Marshal(frameUpdate) //canFrame
		if err != nil {
			fmt.Println("[ERROR] Marshalling JSON: ", err)
		}

		err = connWS.conn.WriteMessage(websocket.TextMessage, jsonData)
		if err != nil {
			fmt.Println("[ERROR] Sending data to client: ", err)
			return
		}
	}
}

func webSocketListener(connWS WtWebSocket) {
	for {
		_, msg, err := connWS.conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				fmt.Println("[INFO] WebSocket disconnected..")
				return
			} else {
				fmt.Println("[INFO] Unexpected disconnection: ", err)
			}
		}

		var mesRes WtWebSocketResponse
		err = json.Unmarshal(msg, &mesRes)
		if err != nil {
			fmt.Println("[ERROR] Unmarshalling from client: ", err)
			continue
		}

		switch mesRes.Type {
		case 1: // GET-MSG-DATA-{id}
			// Lookup the provided CAN Id, get the data, return to client
			u64, err := strconv.ParseUint(mesRes.Data, 10, 32)
			if err != nil {
				fmt.Println("[ERROR] Converting CAN Id input from client")
				return
			}

			canId := uint32(u64)
			canFrame, exists := cfMap[canId]
			if exists {
				jsonData, err := json.Marshal(canFrame)
				if err != nil {
					fmt.Println("[ERROR] Marshalling JSON: ", err)
					return
				}

				err = connWS.conn.WriteMessage(websocket.TextMessage, jsonData)
				if err != nil {
					fmt.Println("[ERROR] Sending data to client: ", err)
					return
				}
			}

		case 2: // Write dataset to file
			saveDataToFile()
		}
	}
}

func saveDataToFile() {
	dir, err := os.Getwd()
    if err != nil {
        fmt.Println("[ERROR] Unable to get current working directory: ", err)
		return
    }

	parentDir := filepath.Dir(dir)

	f, err := os.Create(filepath.Join(parentDir, "wt-can-sniffer-output.json"))
	if err != nil {
		fmt.Println("[ERROR] Creating JSON file: ", err)
		return
	}
	defer f.Close()

	json.NewEncoder(f).Encode(cfMap)

	fmt.Println("[INFO] Written file successfully")

}
