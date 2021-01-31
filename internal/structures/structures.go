package structures

type Status string

const (
	InWork  Status = "in work"
	Stopped Status = "stopped"
)

type Task struct {
	UUID     string `json:"uuid"`
	CameraIP string `json:"camera_ip"`
	ADDR     string `json:"addr"`
	Status   Status `json:"status"`
	Job      string `json:"job"`
}

type State struct {
	ISWork bool
}
