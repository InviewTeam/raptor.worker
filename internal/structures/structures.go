package structures

import "github.com/deepch/vdk/av"

type Status string

const (
	InWork  Status = "start"
	Stopped Status = "stop"
)

type Task struct {
	UUID     string `json:"uuid"`
	CameraIP string `json:"camera_ip"`
	Status   Status `json:"status"`
	Job      []Job  `json:"job"`
}

type State struct {
	ISWork bool
}

type Job struct {
	name    string
	address string
}

type ConfigST struct {
}
type StreamST struct {
	URL    string `json:"url"`
	Codecs []av.CodecData
	Cl     map[string]viewer
}

type viewer struct {
	c chan av.Packet
}
