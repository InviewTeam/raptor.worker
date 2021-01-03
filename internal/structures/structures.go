package structures

type Task struct {
	UUID     string   `json:"uuid"`
	CameraIP string   `json:"camera_ip"`
	Jobs     []string `json:"jobs"`
	Status   string   `json:"status"`
}

type State struct {
	ISWork bool
}
