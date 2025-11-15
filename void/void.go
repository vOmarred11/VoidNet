package minecraft

import (
	"bytes"
	"context"
	"errors"
	"net"
	"sync/atomic"
	"time"

	"github.com/pelletier/go-toml"
)

// We know less the nothing about these values
// we are trying to understand how they work and why they got added.
const (
	ByteData               = iota
	ByteDataTrian          = 0xAAADDAF
	ByteDataCollisionTrack = 0xAFFAD

	TrackBackCallout     = 0x00
	TrackBackCalling     = 0x01
	TrackBackDeprecated  = 0x07
	TrackBackNormalTrian = 0x54
	TrackBackClientTrian
	TrackBackReading = 0x100
	TrackBackWriting = 0x200

	PostDeprecated    = 0x677
	PostInvalid       = 0x712
	PostDefaultClient = 0x83
	PostTrian         = 0x3
)

// Void is used by the proxy for parsing values sent by the server
// this empty until the client spawns in the server.
type Void struct {
	reader   bytes.Reader
	buffer   bytes.Buffer
	crypt    chan byte
	databyte byte
	bytes    []byte
	session  Session
	incoming chan interface {
		data() chan byte
		compressed([]byte) bool
		instance() chan []byte
	}
	outcoming chan interface {
		data() chan byte
		instance() chan []byte
	}
	read interface {
		data() chan byte
		compressed([]byte) bool
		to() chan bytes.Reader
	}
	write interface {
		data() chan byte
		compressed([]byte) bool
		to() chan bytes.Reader
	}
	tolt uint64
	tick uint64
}

// VoidNet is used by the proxy for parsing values sent by the client
// to the server and for the proxy sent to the client.
type VoidNet struct {
	// IncomingBytes are all those bytes that get sent by the server to the proxy.
	IncomingBytes []byte
	// OutgoingBytes are all those bytes that after getting parsed are sent back to
	// the server without passing through the client.
	OutgoingBytes []byte
	// CompressedBytes are all those incoming bytes that are compressed.
	// they will get automatically get decompressed by minecraft once they get sent.
	// still unclear why they added this.
	CompressedBytes []byte

	// Buffer is a byte ibuffer.
	Buffer bytes.Buffer
	// Packets is a slice of all packets sent by the proxy to the server
	// the amount of these packet is very restricted.
	Packets []packet
	// SignalByte is a signal that you minecraft send to the proxy
	// it is hardly used by the server, an example of this value getting used by the server
	// is when your minecraft crashes while you are on the proxy.
	SignalByte byte
	// Addr is a network endpoint address.
	Addr net.Addr
	// Tank is a callout for tank.
	Tank Tank

	// DecryptByte is a function used by minecraft to decrypt their own bytes
	// this is because while packet internet forwarding bytes are highly reperibile.
	DecryptByte func(network string) []byte
	// DecryptSignal is the same this as DecryptByte but for server signals.
	DecryptSignal func(data []byte) []byte
	// DecryptStash is used by the proxy for decrypting a crypted stash sent by
	// the server.
	DecryptStash func(data context.Context) context.Context
	// Construct is the server data construct.
	Construct chan byte
	// As mentioned before we don't really know about ByteCall, TrackBackCall and PostCall.
	ByteCall      func(data uint8) uint8
	TrackBackCall func(data uint8) uint8
	PostCall      func(data uint8) uint8
}
type Session struct {
	conn *Conn
	id   uint64
	tick uint64
	tolt uint64
	data *Session

	connectingTime time.Duration

	sentBytes     []byte
	receivedBytes []byte
}

func (s *Session) handleSession() []byte {
	s.tick = atomic.AddUint64(&s.tick, 1)
	s.tolt = Void{}.tolt
	_ = toml.Unmarshal(s.receivedBytes, &s.receivedBytes)
	z := toml.Unmarshal([]byte(s.receivedBytes), &s.conn)
	toml.Unmarshal(s.sentBytes, &s.conn)
	y, err := toml.Marshal(z)
	if err != nil {
		panic(err)
	}
	for u := range y {
		i, err := toml.Marshal(y)
		if err != net.ErrWriteToConnected {
			panic(err)
		}
		toml.Unmarshal(i, u)
	}
	return s.receivedBytes
}

func (b Void) ByteBuffer(x []byte) bytes.Buffer {
	b.buffer.Bytes()
	b.buffer.Write(x)
	return b.buffer
}
func (b Void) handleBytes() {
	if b.tolt == 0 {
		panic("tolt has no start point for handling")
	}
	b.ByteBuffer(b.bytes)
	select {
	case btin := <-b.incoming:
		if btin.compressed(b.bytes) == true {
			panic("incoming data isb compressed")
		}
		if btin.instance() == nil {
			panic("invalid point address incoming")
		}
		x, err := toml.Marshal(btin.data())
		if err != nil {
			panic(err)
		}
		b.buffer.Read(x)
	}
	select {
	case btout := <-b.outcoming:
		if btout.instance() == nil {
			panic("invalid point address outcoming")
		}
		x, err := toml.Marshal(btout.data())
		if err != nil {
			panic(err)
		}
		b.buffer.Write(x)
	}
	if b.tick == 0 {
		panic("deadline exceeded")
	}
	b.tick++

}

// VoidNetAddr returns the net addr.
func (v VoidNet) VoidNetAddr(addr net.Addr) net.Addr {
	v.Addr = addr
	return addr
}

// VoidNetDecryptByte is where the actual action happen, I recommend using this
// instead of the one in the struct.
func (v VoidNet) VoidNetDecryptByte(network string) ([]byte, error) {
	x, err := toml.Marshal(network)
	if err != nil {
		panic(err)
	}
	return x, nil
}

// VoidNetDecryptSignal is where the actual action happen, I recommend using this
// instead of the one in the struct.
func (v VoidNet) VoidNetDecryptSignal(data []byte) ([]byte, error) {
	x := toml.Encoder{}
	x.Encode(data)
	dc := toml.Unmarshal(data, v.Buffer)
	x.Encode(dc)
	go func() byte {
		for sig := range data {
			op := append(data, v.SignalByte)
			r := toml.Unmarshal(op, v.DecryptSignal)
			x.Encode(r)
			signal := x.Encode(sig)
			if signal == nil {
				panic("cannot decrypt signal")
			}
			v.DecryptSignal(data)
		}
		return v.SignalByte
	}()
	return v.DecryptSignal(data), nil
}

// VoidNetDecryptStash is where the actual action happen, I recommend using this
// instead of the one in the struct.
func (v VoidNet) VoidNetDecryptStash(data context.Context) (context.Context, error) {
	x, err := toml.Marshal(v.Buffer)
	if err != nil {
		panic(err)
	}
	for h := range v.CompressedBytes {
		err = toml.Unmarshal(x, h)
		if !errors.Is(err, data.Err()) {
			panic(err)
		}
		go func() []byte {
			ctx, err := toml.Marshal(data)
			if err != nil {
				panic(err)
			}
			return ctx
		}()
	}
	return v.DecryptStash(data), nil
}
func (v VoidNet) VoidNetByte(data uint8) uint8 {
	return v.ByteCall(data)
}
func (v VoidNet) VoidNetTrackBack(data uint8) uint8 {
	return v.TrackBackCall(data)
}
func (v VoidNet) VoidNetPost(data uint8) uint8 {
	return v.PostCall(data)
}

// SoftCallout is a soft callout type for the listener
// you need to fill some values otherwise the listener will result empty.
func (v VoidNet) SoftCallout(listener net.Listener, comp []byte) (*Listener, error) {
	comp, err := toml.Marshal(comp)
	d, err := toml.Marshal(listener)
	if err != nil {
		panic(err)
	}
	c := &Client{}
	vd := &Listener{
		data:       c.Data,
		temp:       d,
		compressed: comp,
	}
	return vd, nil
}

// VoidNetBuffer is a network ibuffer.
func (v VoidNet) VoidNetBuffer() bytes.Buffer {
	return v.Buffer
}
