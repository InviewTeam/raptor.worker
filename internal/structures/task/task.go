package task

type Task struct {
	UUID     string   `json:"uuid"`
	CameraIP string   `json:"camera_ip"`
	Jobs     []string `json:"jobs"`
	Status   string   `json:"status"`
}
