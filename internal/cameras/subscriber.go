package cameras

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/deepch/vdk/codec/h264parser"
	"io"
	"net/http"
	"time"

	"github.com/deepch/vdk/av"
	"github.com/deepch/vdk/format/rtsp"
	"github.com/pion/webrtc/v3"
	"github.com/pion/webrtc/v3/pkg/media"
	"gitlab.com/inview-team/raptor_team/worker/internal/logger"
)

var (
	outboundVideoTrack *webrtc.TrackLocalStaticSample
)

func WorkWithVideo(url string) {
	var err error
	outboundVideoTrack, err = webrtc.NewTrackLocalStaticSample(webrtc.RTPCodecCapability{
		MimeType: "video/h264",
	}, "pion-rtsp", "pion-rtsp")

	if err != nil {
		logger.Critical.Panic(err.Error())
	}

	go rtspConsumer(url)
}

func rtspConsumer(url string) error {
	annexbNALUStartCode := func() []byte { return []byte{0x00, 0x00, 0x00, 0x01} }
	for {
		session, err := rtsp.Dial(url)
		if err != nil {
			logger.Critical.Panic(err.Error())
		}
		session.RtpKeepAliveTimeout = 10 * time.Second
		codecs, err := session.Streams()
		if err != nil {
			logger.Critical.Panic(err.Error())
		}

		for i, t := range codecs {
			logger.Info.Println("Stream", i, "is of type", t.Type().String())
		}

		if codecs[0].Type() != av.H264 {
			logger.Critical.Panic("RTSP feed must begin with a H264 codec")
		}

		if len(codecs) != 1 {
			logger.Info.Println("Ignoring all but the first stream.")
		}

		var previousTime time.Duration

		logger.Info.Println("Start streaming ", url)

		for {
			pkt, err := session.ReadPacket()
			if err != nil {
				break
			}

			if pkt.Idx != 0 {
				continue
			}

			pkt.Data = pkt.Data[4:]

			// For every key-frame pre-pend the SPS and PPS
			if pkt.IsKeyFrame {
				pkt.Data = append(annexbNALUStartCode(), pkt.Data...)
				pkt.Data = append(codecs[0].(h264parser.CodecData).PPS(), pkt.Data...)
				pkt.Data = append(annexbNALUStartCode(), pkt.Data...)
				pkt.Data = append(codecs[0].(h264parser.CodecData).SPS(), pkt.Data...)
				pkt.Data = append(annexbNALUStartCode(), pkt.Data...)
			}

			bufferDuration := pkt.Time - previousTime
			previousTime = pkt.Time
			if err = outboundVideoTrack.WriteSample(media.Sample{Data: pkt.Data, Duration: bufferDuration}); err != nil && err != io.ErrClosedPipe {
				panic(err)
			}
		}

		logger.Info.Println("Stream stopped")
		if err = session.Close(); err != nil {
			logger.Info.Println("session Close error", err)
		}

		time.Sleep(5 * time.Second)
	}
}

func DoSignaling(w http.ResponseWriter, r *http.Request) {
	peerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		logger.Critical.Panic(err.Error())
	}

	if _, err = peerConnection.AddTrack(outboundVideoTrack); err != nil {
		logger.Critical.Panic(err.Error())
	}

	var offer webrtc.SessionDescription
	if err = json.NewDecoder(r.Body).Decode(&offer); err != nil {
		logger.Critical.Panic(err.Error())
	}

	if err = peerConnection.SetRemoteDescription(offer); err != nil {
		logger.Critical.Panic(err.Error())
	}

	gatherCompletePromise := webrtc.GatheringCompletePromise(peerConnection)

	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		panic(err)
	} else if err = peerConnection.SetLocalDescription(answer); err != nil {
		panic(err)
	}

	<-gatherCompletePromise

	response, err := json.Marshal(*peerConnection.LocalDescription())
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(response); err != nil {
		panic(err)
	}
}

func NewPeerConnection() {
	peerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		panic(err)
	}

	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		fmt.Printf("Connection State has changed %s \n", connectionState.String())
	})

	offer, err := peerConnection.CreateOffer(nil)
	if err != nil {
		panic(err)
	}

	gatherComplete := webrtc.GatheringCompletePromise(peerConnection)

	if err = peerConnection.SetLocalDescription(offer); err != nil {
		panic(err)
	}

	offerJSON, err := json.Marshal(*peerConnection.LocalDescription())
	if err != nil {
		panic(err)
	}

	resp, err := http.Post("http://0.0.0.0:8080/offer", "application/json", bytes.NewReader(offerJSON))
	if err != nil {
		panic(err)
	}

	resp.Close = true

	var answer webrtc.SessionDescription

	if err = json.NewDecoder(resp.Body).Decode(&answer); err != nil {
		logger.Critical.Panic(err.Error())
	}

	resp.Body.Close()

	<-gatherComplete

}
