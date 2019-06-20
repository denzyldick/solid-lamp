package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
	"net/http"
	"strconv"
)

var addr = flag.String("addr", ":4444", "http service address")
// We'll need to define an Upgrader
// this will require a Read and Write buffer size
var upgrader = websocket.Upgrader{

}
func main() {
	flag.Parse()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {

		serve( w, r)
	})
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Print("ListenAndServe: ", err)
	}
}

var clients []Client
// serve handles websocket requests from the peer.
func serve( w http.ResponseWriter, r *http.Request) {
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

		fmt.Println("new message received")
		// print out that message for clarity
		payload  := Payload{}
		err = json.Unmarshal(p,&payload)
		fmt.Println(payload.Type)
		if err!= nil{
			log.Print(err)
		}


		if payload.Type == "SDP"{
			client := Client{
				SDP:payload.SDP,
				DispayName: strconv.Itoa(rand.Int()),
				Conn: *conn,
			}
			clients = append(clients, client)

		}

		if payload.Type == "NEW_CANDIDATE" {

			for c := range clients {

					//bytes, err := json.Marshal(clients[c])

					if err != nil {
						log.Print(err)
					}
					err = conn.WriteJSON(clients[c])

					if err != nil {
						log.Print(err)
					} else {
						fmt.Print("Mesage send")
					}

			}
		}
		}
		//if err := conn.WriteMessage(messageType, p); err != nil {
		//	log.Println(err)
		//	return
		//}


}

//func findClients(start *Client)[]*Client{
//	return &Client{}
//}
type Payload struct{
	Type string `json:"type"`
	Value string `json:"value"`
	SDP SDP `json:"sdp"`
}

type SDP struct{
	Type string `json:"type"`
	Sdp string	`json:"sdp"`
}
type Client struct {
	SDP SDP
	DispayName string
	Conn websocket.Conn
}


type Group struct {
	Members []Client
}
