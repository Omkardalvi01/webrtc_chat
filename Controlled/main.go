package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/pion/webrtc/v3"
)

func main() {
	peerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	})
	if err != nil {
		log.Fatal("Error while creating peerconnection ", err)
	}

	peerConnection.OnDataChannel(func(dc *webrtc.DataChannel) {
		fmt.Printf("New DataChannel %s\n", dc.Label())

		dc.OnOpen(func() {
			fmt.Println("Connected to peer. Type messages:")

			go func() {
				for {
					msg, _ := bufio.NewReader(os.Stdin).ReadString('\n')
					dc.SendText(msg)
				}
			}()
		})

		dc.OnMessage(func(msg webrtc.DataChannelMessage) {
			fmt.Printf("Peer:%v", string(msg.Data))
		})
	})

	fmt.Println("Put SDP for remote server below")
	offer, err := bufio.NewReader(os.Stdin).ReadString('#')
	if err != nil {
		log.Fatal("Error at reading input ", err)
	}

	offerstr := strings.TrimSpace(offer)
	offerSDP := webrtc.SessionDescription{
		SDP:  offerstr,
		Type: webrtc.SDPTypeOffer,
	}

	fmt.Println("Offer recieved:")
	fmt.Println(offerSDP.Type)

	err = peerConnection.SetRemoteDescription(offerSDP)
	if err != nil {
		log.Fatal("Error while setting remote description ", err)
	}

	ans, err := peerConnection.CreateAnswer(nil)
	if err != nil  {
		log.Fatal("Error while creating answer for new connection ", err)
	}

	if err := peerConnection.SetLocalDescription(ans); err != nil {
		log.Fatal("Error while setting local desciption ", err)
	}
	gathercomplete := webrtc.GatheringCompletePromise(peerConnection)
	<-gathercomplete
	fmt.Printf("Local SDP:\n%v", peerConnection.LocalDescription().SDP)

	select {}
}
