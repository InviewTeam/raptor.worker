package cameras

import (
	"github.com/nareix/joy4/format/rtsp"
	"gitlab.com/inview-team/raptor_team/worker/internal/logger"
)

func checkConnectionToCamera(url string) bool {
	logger.Info.Printf("Check connection to", url)
	rtsp.DebugRtsp = true
	_, err := rtsp.Dial(url)
	if err != nil {
		logger.Error.Printf("Connection to this camera unavailable", url)
		return false
	}
	return true
}
