package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/pion/webrtc/v3"
)

func main(){
	webconfig := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}
	
	peerConnection , err := webrtc.NewPeerConnection(webconfig)
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
		fmt.Printf("Peer: %s\n", string(msg.Data))
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

	fmt.Printf("Local SDP:\n%v", peerConnection.LocalDescription().SDP)

	fmt.Println("Enter the remote SDP address below")
	ansSDP , err := bufio.NewReader(os.Stdin).ReadString('#')
	if err != nil {
		log.Fatal("Error while reading ",err)
	}

	ansSDP_str := strings.TrimSpace(string(ansSDP))
	answer := webrtc.SessionDescription{
		SDP: ansSDP_str,
		Type: webrtc.SDPTypeAnswer,
	}

	if err := peerConnection.SetRemoteDescription(answer); err != nil{
		log.Fatal("Error while setting remote description ",err)
	}

	select{}
}