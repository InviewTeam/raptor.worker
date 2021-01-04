package structures

type Task struct {
	UUID     string `json:"uuid"`
	CameraIP string `json:"camera_ip"`
	ADDR     string `json:"addr"`
	Status   string `json:"status"`
}

type State struct {
	ISWork bool
}
