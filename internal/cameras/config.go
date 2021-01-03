package cameras

import (
	"crypto/rand"
	"fmt"
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

func (element *ConfigST) coGe() []av.CodecData {
	return element.Stream.Codecs
}

func (element *ConfigST) coAd(codecs []av.CodecData) {
	element.Stream.Codecs = codecs
}

func (element *ConfigST) clAd() (string, chan av.Packet) {
	cuuid := pseudoUUID()
	ch := make(chan av.Packet, 100)
	element.Stream.Cl[cuuid] = viwer{c: ch}
	return cuuid, ch
}

func (element *ConfigST) clDe(cuuid string) {
	delete(element.Stream.Cl, cuuid)
}

func pseudoUUID() (uuid string) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	uuid = fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return
}
