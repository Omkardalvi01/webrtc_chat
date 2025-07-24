package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)


var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
        return true 
    },
}

var Peer = make(map[string]chan string)

func exchangeSDP(conn *websocket.Conn){
	defer conn.Close()

	_ , id , err := conn.ReadMessage()
	if err != nil{
		fmt.Print("Error while reading id",err)
	}	

	conn_id := string(id)
	fmt.Println("id: ",conn_id) 
	comms_channel := Peer[conn_id]
	
	if comms_channel == nil{
		fmt.Println("Creating new chat")
		comms_channel = make(chan string)
		Peer[conn_id] = comms_channel
	}

	
	select{
	case offer := <-comms_channel:
		err := conn.WriteMessage(websocket.TextMessage, []byte(offer))
		if err != nil{
			fmt.Print("Error forwarding offer ", err)
			return 
		}
		fmt.Println("Directing offer to peerB",string(offer))

		_ , answer , err := conn.ReadMessage()
		if err != nil{
			fmt.Print("Error reading answer ", err)
			return 
		}
		fmt.Println("Sending answer to peerA", string(answer))
		comms_channel <- string(answer)
	
	default:
		_ , offer, err := conn.ReadMessage()
		if err != nil{
			fmt.Print("Error reading offer ", err)
			return 
		}
		fmt.Println("Recieved Offer ", string(offer))
		comms_channel <- string(offer)

		answer := <-comms_channel
		err = conn.WriteMessage(websocket.TextMessage, []byte(answer)) 
		if err != nil{
			fmt.Print("Error at writing answer ", err)
			return 
		}
		fmt.Println("Recieved answer from peerB", string(answer))
	}
}

func wshandler(w http.ResponseWriter, r *http.Request){
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil{
		fmt.Println("Couldnt upgrade to websocket ",err)
		return
	}

	exchangeSDP(conn)
}

func main(){

	http.HandleFunc("/", wshandler)
	log.Fatal(http.ListenAndServe(":10000",nil))
}