package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":4444", "http service address")

// We'll need to define an Upgrader
// this will require a Read and Write buffer size
var upgrader = websocket.Upgrader{}

func main() {
	flag.Parse()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {

		serve(w, r)
	})
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Print("ListenAndServe: ", err)
	}
}

var clients []Client

// serve handles websocket requests from the peer.
func serve(w http.ResponseWriter, r *http.Request) {
	/// Allow cross origin for now.
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, w.Header())
	if err != nil {
		log.Println(err)
		return
	}

	reader(conn)
}

// define a reader which will listen for
// new messages being sent to our WebSocket
// endpoint
func reader(conn *websocket.Conn) {
	for {
		// read in a message
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		// print out that message for clarity
		payload := Payload{}
		err = json.Unmarshal(p, &payload)
		fmt.Println(payload.Type)
		if err != nil {
			log.Print(err)
		}
		client := Client{
			SDP:        payload.SDP,
			DispayName: payload.Id,
			Conn:       *conn,
		}
		if payload.Type == "SDP" {

		}

		if findClients(&client) == false {
			clients = append(clients, client)
		}
		for c := range clients {
			err = conn.WriteJSON(clients[c])

			if err != nil {
				log.Print(err)
			} else {
				fmt.Println(payload)
			}

		}

		fmt.Println(len(clients), "clients")
	}
}

func findClients(start *Client) bool {

	for i, _ := range clients {
		return clients[i].DispayName == start.DispayName
	}
	return false

}

type Payload struct {
	Id string `json:"id"`
	Type  string `json:"type"`
	Value string `json:"value"`
	SDP   SDP    `json:"sdp"`
}

type SDP struct {
	Type string `json:"type"`
	Sdp  string `json:"sdp"`
}
type Client struct {
	SDP        SDP
	DispayName string
	Conn       websocket.Conn
}

type Group struct {
	Members []Client
}
