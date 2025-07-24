package utils

import (
	"context"
	"time"
	"github.com/gorilla/websocket"
)

func Send_uid(uid string) (*websocket.Conn, error){
	
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(5 * time.Second))
	defer cancel()
	
	conn, _ , err := websocket.DefaultDialer.DialContext(ctx, "ws://localhost:8080/" , nil)
	if err != nil{
		return nil , err
	}

	err = conn.WriteMessage(websocket.TextMessage, []byte(uid))
	if err != nil{
		return nil, err
	}

	return conn , nil
}

func Send_and_recieve(conn *websocket.Conn, data string) (string , error) {

	err := conn.WriteMessage(websocket.TextMessage, []byte(data))
	if err != nil{
		return "", err
	}

	_ , resp , err := conn.ReadMessage()
	if err != nil{
		return "", err
	}
	return string(resp), nil
}

func Send(conn *websocket.Conn, data string) error {

	err := conn.WriteMessage(websocket.TextMessage, []byte(data))
	if err != nil{
		return err
	}
	return nil
}

func Recieve(conn *websocket.Conn) (string , error) {
	_, resp, err := conn.ReadMessage()
	if err != nil {
		return "", nil
	}
	return string(resp) , nil
}