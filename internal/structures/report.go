package structures

type Report struct {
	UUID     string   `json:"uuid"`
	Analyzer string   `json:"analyzer"`
	Results  []Result `json:"results"`
}

type Result struct {
	Time string `json:"time"`
	Info string `json:"info"`
}
