package cameras

import (
	"encoding/json"
	"github.com/deepch/vdk/av"
	"github.com/deepch/vdk/codec/h264parser"
	"github.com/deepch/vdk/format/rtsp"
	"github.com/pion/webrtc/v3"
	"github.com/pion/webrtc/v3/pkg/media"
	"gitlab.com/inview-team/raptor_team/worker/internal/logger"
	"io"
	"net/http"
	"time"
)

var (
	outboundVideoTrack *webrtc.TrackLocalStaticSample
)

func WorkWithVideo(url string, addr string, ch <-chan bool) {
	var err error
	outboundVideoTrack, err = webrtc.NewTrackLocalStaticSample(webrtc.RTPCodecCapability{
		MimeType: "video/h264",
	}, "pion-rtsp", "pion-rtsp")
	if err != nil {
		logger.Critical.Panic(err.Error())
	}

	go rtspConsumer(url, ch)
	http.Handle("/", http.FileServer(http.Dir("C://Users/sitadm/GolandProjects/worker/internal/cameras/static")))
	http.HandleFunc("/doSignaling", doSignaling)

	logger.Info.Println("Starting stream")
	logger.Critical.Panic(http.ListenAndServe(addr, nil))
}

func rtspConsumer(url string, ch <-chan bool) {
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

		for {
			select {
			case msg := <-ch:
				if msg {
					return
				}
			default:
				continue
			}

			pkt, err := session.ReadPacket()
			if err != nil {
				break
			}

			if pkt.Idx != 0 {
				continue
			}

			pkt.Data = pkt.Data[4:]

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
				logger.Critical.Panic("RTSP feed must begin with a H264 codec")
			}
		}

		if err = session.Close(); err != nil {
			logger.Info.Println("session Close error", err)
		}

		time.Sleep(5 * time.Second)
	}
}

func doSignaling(w http.ResponseWriter, r *http.Request) {
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
