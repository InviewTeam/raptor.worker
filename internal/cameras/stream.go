package cameras

import (
	"github.com/deepch/vdk/format/rtsp"
	"gitlab.com/inview-team/raptor_team/worker/internal/logger"
	"time"
)

func ServeStream(url string) {
	for {
		logger.Info.Println("Connect to ", url)
		rtsp.DebugRtsp = true
		session, err := rtsp.Dial(url)
		if err != nil {
			logger.Error.Println(err.Error())
			time.Sleep(5 + time.Second)
			continue
		}
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
