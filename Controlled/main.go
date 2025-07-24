package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"github.com/p2p_webrtc_chat/Utils"
	"github.com/pion/webrtc/v3"
)

func main() {
	os_reader := bufio.NewReader(os.Stdin)
	peer_with_escape , _ := os_reader.ReadString('\n')
	peer :=  peer_with_escape[:len(peer_with_escape)-1]

	peerConnection, err := webrtc.NewPeerConnection(utils.Webconfig)
	if err != nil {
		log.Fatal("Error while creating peerconnection ", err)
	}

	peerConnection.OnDataChannel(func(dc *webrtc.DataChannel) {
		fmt.Printf("New DataChannel %s\n", dc.Label())

		dc.OnOpen(func() {
			fmt.Println("Connected to peer. Type messages:")

			go func() {
				for {
					msg, _ := os_reader.ReadString('\n')
					dc.SendText(msg)
				}
			}()
		})

		dc.OnMessage(func(msg webrtc.DataChannelMessage) {
			fmt.Printf("%s:%v",peer, string(msg.Data))
		})
	})


	fmt.Println("Put the id below:")
	uid, err := os_reader.ReadString('\n')
	if err != nil {
		log.Fatal("Error at reading input ", err)
	}


	conn , err := utils.Send_uid(uid[:len(uid)-1])
	if err != nil{
		log.Fatal("Error with creating connection", err)
	}
	

	offer, err := utils.Recieve(conn)
	if err != nil{
		log.Fatal("Error with recieving offer", err)
	}

	offerSDP := webrtc.SessionDescription{
		SDP:  offer,
		Type: webrtc.SDPTypeOffer,
	}

	fmt.Println("Offer recieved:")
	fmt.Println(offerSDP.SDP)

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
	
	err = utils.Send(conn, peerConnection.LocalDescription().SDP)
	if err != nil{
		log.Fatal("Error while sending answer", err)
	}
	fmt.Printf("Local SDP:\n%v", peerConnection.LocalDescription().SDP)

	conn.Close()
	select {}
}
