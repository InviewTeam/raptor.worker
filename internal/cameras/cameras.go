package cameras

import (
	"bytes"
	"encoding/json"
	"github.com/deepch/vdk/format/rtspv2"
	"github.com/pion/webrtc/v3"
	"gitlab.com/inview-team/raptor_team/worker/internal/logger"
	"net/http"
	"time"
)

func RTSPLoop(url string, uuid string, done chan struct{}) {
	logger.Info.Println("Start streaming: ", uuid)
	for {
		select {
		case <-done:
			return
		default:
			err := GetImagesFromRTSP(url, done)
			if err != nil {
				logger.Error.Panic(err.Error())
			}
		}
	}
}

func GetImagesFromRTSP(cameraUrl string, done chan struct{}) error {

	RTSPClient, err := rtspv2.Dial(rtspv2.RTSPClientOptions{URL: cameraUrl, DisableAudio: true, DialTimeout: 3 * time.Second, ReadWriteTimeout: 3 * time.Second, Debug: false})
	if err != nil {
		return err
	}

	print(RTSPClient.CodecData[0].Type().String())

	defer RTSPClient.Close()

	videoTrack := newPeerConnection("http://0.0.0.0:8080/offer")
	for {
		select {
		case <-done:
			break
		case packet := <-RTSPClient.OutgoingPacketQueue:
			if _, err = videoTrack.Write(packet.Data); err != nil {
				logger.Critical.Panic(err.Error())
			}
		}
	}
	logger.Info.Println("Stop streaming")
	return nil
}

func newPeerConnection(addr string) *webrtc.TrackLocalStaticRTP {
	peerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	})

	if err != nil {
		logger.Error.Panic(err.Error())
	}

	videoTrack, err := webrtc.NewTrackLocalStaticRTP(webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeH264}, "video", "pion")

	if err != nil {
		logger.Error.Panic(err.Error())
	}

	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		logger.Info.Printf("Connection State has changed %s \n", connectionState.String())
	})

	offer, err := peerConnection.CreateOffer(nil)
	if err != nil {
		logger.Error.Panic(err.Error())
	}

	if err = peerConnection.SetLocalDescription(offer); err != nil {
		logger.Error.Panic(err.Error())
	}

	offerJSON, err := json.Marshal(*peerConnection.LocalDescription())
	if err != nil {
		logger.Error.Panic(err.Error())
	}

	resp, err := http.Post(addr, "application/json", bytes.NewReader(offerJSON))
	if err != nil {
		logger.Error.Panic(err.Error())
	}
	resp.Close = true

	gatherComplete := webrtc.GatheringCompletePromise(peerConnection)

	var answer webrtc.SessionDescription

	if err = json.NewDecoder(resp.Body).Decode(&answer); err != nil {
		logger.Critical.Panic(err.Error())
	}

	<-gatherComplete

	if err = peerConnection.SetRemoteDescription(answer); err != nil {
		panic(err)
	}

	resp.Body.Close()

	return videoTrack
}
