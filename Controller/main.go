package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"github.com/google/uuid"
	"github.com/p2p_webrtc_chat/Utils"
	"github.com/pion/webrtc/v3"
)

func main(){
	os_reader := bufio.NewReader(os.Stdin)
	peer_with_escape , _ := os_reader.ReadString('\n')
	peer :=  peer_with_escape[:len(peer_with_escape)-1]

	peerConnection , err := webrtc.NewPeerConnection(utils.Webconfig)
	if err != nil{
		log.Fatal("Error while creating new connection",err)
	}
	defer peerConnection.GracefulClose()

	dc , err := peerConnection.CreateDataChannel("chat",nil)
	if err != nil {
		log.Fatal("Error while creating data channel", err)
	}
	dc.OnOpen(func() {
		fmt.Println("DataChannel is open")
		
		go func(){
			for{
				msg , _ := bufio.NewReader(os.Stdin).ReadString('\n')
				dc.SendText(msg)
			}
		}()
		
	})
	dc.OnMessage(func(msg webrtc.DataChannelMessage) {
		fmt.Printf("%s: %s", peer ,string(msg.Data))
	})

	offer , err := peerConnection.CreateOffer(nil)
	if err != nil{
		log.Fatal("Error while creating offer ",err)
	}

	err = peerConnection.SetLocalDescription(offer)
	if err != nil{
		log.Fatal("Error at setting sdp to local peer",err)
	}
	gathercomplete := webrtc.GatheringCompletePromise(peerConnection)
	<-gathercomplete

	uid := uuid.New().URN()
	fmt.Println("Copy paste to peer:")
	fmt.Println(uid)


	conn , err := utils.Send_uid(uid)
	if err != nil{
		log.Fatal("Error at creating connection to server",err)
	}

	ans , err := utils.Send_and_recieve(conn , peerConnection.LocalDescription().SDP)
	if err != nil{
		log.Fatal("Error at send and recieve to controlled",err)
	}

	answer := webrtc.SessionDescription{
		SDP: ans,
		Type: webrtc.SDPTypeAnswer,
	}

	if err := peerConnection.SetRemoteDescription(answer); err != nil{
		log.Fatal("Error while setting remote description ",err)
	}

	conn.Close()
	select{}
}