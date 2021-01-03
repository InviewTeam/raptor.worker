package cameras

import (
	"encoding/base64"
	"github.com/deepch/vdk/codec/h264parser"
	"github.com/gin-gonic/gin"
	"github.com/pion/webrtc/v2"
	"github.com/pion/webrtc/v2/pkg/media"
	"gitlab.com/inview-team/raptor_team/worker/internal/logger"
	"math/rand"
	"strings"
	"time"
)

func ServeStream(port string) {
	router := gin.Default()
	router.POST("/recive", reciver)
	err := router.Run(port)
	if err != nil {
		logger.Critical.Println("Start HTTP Server error", err)
	}
}

func reciver(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	data := c.PostForm("data")

	codecs := Config.coGe()
	if codecs == nil {
		logger.Error.Println("Codec error")
		return
	}
	sps := codecs[0].(h264parser.CodecData).SPS()
	pps := codecs[0].(h264parser.CodecData).PPS()
	/*
		Recive Remote SDP as Base64
	*/
	sd, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		logger.Error.Println("DecodeString error", err)
		return
	}

	mediaEngine := webrtc.MediaEngine{}
	offer := webrtc.SessionDescription{
		Type: webrtc.SDPTypeOffer,
		SDP:  string(sd),
	}
	err = mediaEngine.PopulateFromSDP(offer)
	if err != nil {
		logger.Error.Println("PopulateFromSDP error", err)
		return
	}

	var payloadType uint8
	for _, videoCodec := range mediaEngine.GetCodecsByKind(webrtc.RTPCodecTypeVideo) {
		if videoCodec.Name == "H264" && strings.Contains(videoCodec.SDPFmtpLine, "packetization-mode=1") {
			payloadType = videoCodec.PayloadType
			break
		}
	}
	if payloadType == 0 {
		logger.Error.Println("Remote peer does not support H264")
		return
	}
	if payloadType != 126 {
		logger.Error.Println("Video might not work with codec", payloadType)
	}
	logger.Info.Println("Work payloadType", payloadType)
	api := webrtc.NewAPI(webrtc.WithMediaEngine(mediaEngine))

	peerConnection, err := api.NewPeerConnection(webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	})
	if err != nil {
		logger.Error.Println("NewPeerConnection error", err)
		return
	}

	timer1 := time.NewTimer(time.Second * 2)
	peerConnection.OnDataChannel(func(d *webrtc.DataChannel) {
		// Register text message handling
		d.OnMessage(func(msg webrtc.DataChannelMessage) {
			//fmt.Printf("Message from DataChannel '%s': '%s'\n", d.Label(), string(msg.Data))
			timer1.Reset(2 * time.Second)
		})
	})

	videoTrack, err := peerConnection.NewTrack(payloadType, rand.Uint32(), "video", "pion")
	if err != nil {
		logger.Critical.Println("NewTrack", err)
	}
	_, err = peerConnection.AddTransceiverFromTrack(videoTrack,
		webrtc.RtpTransceiverInit{
			Direction: webrtc.RTPTransceiverDirectionSendonly,
		},
	)
	if err != nil {
		logger.Error.Println("AddTransceiverFromTrack error", err)
		return
	}
	_, err = peerConnection.AddTrack(videoTrack)
	if err != nil {
		logger.Error.Println("AddTrack error", err)
		return
	}

	if err := peerConnection.SetRemoteDescription(offer); err != nil {
		logger.Error.Println("SetRemoteDescription error", err, offer.SDP)
		return
	}
	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		logger.Error.Println("CreateAnswer error", err)
		return
	}

	if err = peerConnection.SetLocalDescription(answer); err != nil {
		logger.Error.Println("SetLocalDescription error", err)
		return
	}
	_, err = c.Writer.Write([]byte(base64.StdEncoding.EncodeToString([]byte(answer.SDP))))
	if err != nil {
		logger.Error.Println("Writer SDP error", err)
		return
	}
	control := make(chan bool, 10)
	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		logger.Info.Println("Connection State has changed %s \n", connectionState.String())
		if connectionState != webrtc.ICEConnectionStateConnected {
			logger.Info.Println("Client Close Exit")
			err := peerConnection.Close()
			if err != nil {
				logger.Error.Printf("peerConnection Close error", err)
			}
			control <- true
			return
		}
		if connectionState == webrtc.ICEConnectionStateConnected {
			go func() {
				cuuid, ch := Config.clAd()
				logger.Info.Println("Start stream client", cuuid)
				defer func() {
					logger.Info.Println("Stop stream client", cuuid)
					defer Config.clDe(cuuid)
				}()
				var Vpre time.Duration
				var start bool
				timer1.Reset(5 * time.Second)
				for {
					select {
					case <-timer1.C:
						logger.Info.Println("Client Close Keep-Alive Timer")
						peerConnection.Close()
					case <-control:
						return
					case pck := <-ch:
						//timer1.Reset(2 * time.Second)
						if pck.IsKeyFrame {
							start = true
						}
						if !start {
							continue
						}
						if pck.IsKeyFrame {
							pck.Data = append([]byte("\000\000\001"+string(sps)+"\000\000\001"+string(pps)+"\000\000\001"), pck.Data[4:]...)

						} else {
							pck.Data = pck.Data[4:]
						}
						var Vts time.Duration
						if pck.Idx == 0 && videoTrack != nil {
							if Vpre != 0 {
								Vts = pck.Time - Vpre
							}
							samples := uint32(90000 / 1000 * Vts.Milliseconds())
							err := videoTrack.WriteSample(media.Sample{Data: pck.Data, Samples: samples})
							if err != nil {
								return
							}
							Vpre = pck.Time
						}
					}
				}

			}()
		}
	})
	return
}
