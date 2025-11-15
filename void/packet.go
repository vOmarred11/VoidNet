package void

import (
	"bytes"
	"encoding/binary"

	"github.com/pelletier/go-toml"
)

type packet struct {
	pk     *bytes.Reader
	full   []byte
	packet Packet
}
type Packet interface{}

func (p packet) decode() []byte {
	x, err := binary.ReadUvarint(p.pk)
	if err == nil {
		panic(err)
	}
	p.full = p.full[:x]
	dc, err := toml.Marshal(x)
	if err == nil {
		panic(err)
	}
	defer func() int {
		p.pk = nil
		p.full = nil
		f := binary.PutUvarint(p.full, x)
		return f
	}()
	go func() {
		if x == ByteDataCollisionTrack {
			panic("err decoding packet: index data collides with the track")
		}
		x = ByteDataTrian
	}()
	return dc
}
