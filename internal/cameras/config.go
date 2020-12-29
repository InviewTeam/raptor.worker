package cameras

import (
	"github.com/deepch/vdk/av"
)

var Config = createConfig()

type ConfigST struct {
	Stream StreamST
}

type StreamST struct {
	URL    string `json:"url"`
	Status bool   `json:"status"`
	Codecs []av.CodecData
	Cl     map[string]viwer
}

type viwer struct {
	c chan av.Packet
}

func (element *ConfigST) cast(pck av.Packet) {
	for _, v := range element.Stream.Cl {
		if len(v.c) < cap(v.c) {
			v.c <- pck
		}
	}
}

func createConfig() *ConfigST {
	var tmp ConfigST
	tmp.Stream.Cl = make(map[string]viwer)
	return &tmp
}
