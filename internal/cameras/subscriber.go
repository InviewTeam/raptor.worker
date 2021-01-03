package cameras

import (
	"github.com/deepch/vdk/format/rtsp"
	"gitlab.com/inview-team/raptor_team/worker/internal/logger"
	"time"
)

func SubsribeStream(url string) {
	for {
		logger.Info.Println("Connect to ", url)
		rtsp.DebugRtsp = true
		session, err := rtsp.Dial(url)
		if err != nil {
			logger.Error.Println(err.Error())
			time.Sleep(5 + time.Second)
			continue
		}
		session.RtpKeepAliveTimeout = 10 * time.Second
		if err != nil {
			logger.Error.Println(err.Error())
			time.Sleep(5 * time.Second)
			continue
		}
		codec, err := session.Streams()
		if err != nil {
			logger.Error.Println(err.Error())
			time.Sleep(5 * time.Second)
			continue
		}

		Config.coAd(codec)

		for {
			pkt, err := session.ReadPacket()
			if err != nil {
				logger.Error.Println(err.Error())
				time.Sleep(5 + time.Second)
				continue
			}
			Config.cast(pkt)
		}
		err = session.Close()
		if err != nil {
			logger.Error.Println(err.Error())
		}
		logger.Info.Println("reconnect wait 5s")
		time.Sleep(5 * time.Second)
	}
}
