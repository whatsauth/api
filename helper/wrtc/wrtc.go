package wrtc

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/websocket/v2"
	"github.com/pion/webrtc/v3"
)

var peers = make(map[*websocket.Conn]*webrtc.PeerConnection)

func RunWebRTCSocket(c *websocket.Conn) {
	defer func() {
		if peerConnection, ok := peers[c]; ok {
			peerConnection.Close()
			delete(peers, c)
		}
		c.Close()
	}()

	config := webrtc.Configuration{}
	peerConnection, err := webrtc.NewPeerConnection(config)
	if err != nil {
		fmt.Println("Failed to create peer connection:", err)
		return
	}

	peers[c] = peerConnection

	peerConnection.OnICECandidate(func(candidate *webrtc.ICECandidate) {
		if candidate == nil {
			return
		}
		candidateJSON, _ := json.Marshal(candidate.ToJSON())
		c.WriteMessage(websocket.TextMessage, candidateJSON)
	})

	peerConnection.OnTrack(func(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		// Handle incoming tracks
	})

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			fmt.Println("Error reading message:", err)
			return
		}

		var signal map[string]interface{}
		err = json.Unmarshal(message, &signal)
		if err != nil {
			fmt.Println("Error unmarshalling message:", err)
			continue
		}

		handleSignal(peerConnection, signal, c)
	}
}

func handleSignal(peerConnection *webrtc.PeerConnection, signal map[string]interface{}, c *websocket.Conn) {
	if sdp, ok := signal["sdp"]; ok {
		sdpMap := sdp.(map[string]interface{})
		sdpType := sdpMap["type"].(string)
		sdpContent := sdpMap["sdp"].(string)

		session := webrtc.SessionDescription{
			SDP:  sdpContent,
			Type: webrtc.NewSDPType(sdpType),
		}

		if session.Type == webrtc.SDPTypeOffer {
			if err := peerConnection.SetRemoteDescription(session); err != nil {
				fmt.Println("Error setting remote description:", err)
				return
			}

			answer, err := peerConnection.CreateAnswer(nil)
			if err != nil {
				fmt.Println("Error creating answer:", err)
				return
			}

			if err := peerConnection.SetLocalDescription(answer); err != nil {
				fmt.Println("Error setting local description:", err)
				return
			}

			answerJSON, _ := json.Marshal(answer)
			c.WriteMessage(websocket.TextMessage, answerJSON)
		} else if session.Type == webrtc.SDPTypeAnswer {
			if err := peerConnection.SetRemoteDescription(session); err != nil {
				fmt.Println("Error setting remote description:", err)
				return
			}
		}
	} else if candidate, ok := signal["candidate"]; ok {
		candidateMap := candidate.(map[string]interface{})
		candidateString := candidateMap["candidate"].(string)
		sdpMid := candidateMap["sdpMid"].(string)
		sdpMLineIndex := uint16(candidateMap["sdpMLineIndex"].(float64))

		iceCandidate := webrtc.ICECandidateInit{
			Candidate:     candidateString,
			SDPMid:        &sdpMid,
			SDPMLineIndex: &sdpMLineIndex,
		}

		if err := peerConnection.AddICECandidate(iceCandidate); err != nil {
			fmt.Println("Error adding ICE candidate:", err)
			return
		}
	}
}
