package structures

type Task struct {
	UUID     string `json:"uuid"`
	CameraIP string `json:"camera_ip"`
	Jobs     []Job  `json:"jobs"`
}

type Job struct {
	Title string `json:"title"`
}
