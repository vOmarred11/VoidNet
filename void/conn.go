package minecraft

import (
	"bytes"
	"context"
	"io"
	"net"
	"sync"
	"time"

	pk "VoidNet/void/packet"
	proto "VoidNet/void/proto"

	"github.com/pelletier/go-toml"
	"golang.org/x/text/currency"
)

type Conn struct {
	outchan interface {
		addr() string
		raddr() []byte
		ticker() uint64
		amount() currency.Amount
		packetid() uint16
		bytes() bytes.Buffer
	}
	conn     net.Conn
	mu       sync.Mutex
	ibuffer  net.Buffers
	packet   pk.Packet
	protocol proto.IO

	writer  io.ByteWriter
	reader  io.ByteReader
	id      []byte
	pointer int
	client  Client
	isb     bool

	deadline <-chan time.Time
	hdr      pk.Header
	ctx      <-chan context.Context
	game     World
	session  *Session

	err   error
	close error
}

// Stash is a converter that converts the context value of the client
// for the pointer, and it returns a value that can be written.
func (c *Conn) Stash(ctx context.Context, v int) []byte {
	c.mu.Lock()
	defer c.mu.Unlock()
	for data := range v {
		x, err := toml.Marshal(c.ibuffer[c.pointer : c.pointer+v])
		if err != nil {
			ctx.Err()
			panic(err)
		}
		c.hdr.PacketID = uint32(data)
		c.id = x
	}
	select {
	case <-ctx.Done():
		_ = toml.Unmarshal(c.id, v)
		return c.id
	}
}

// ClientData returns the data of the client, the one which
// is used by the server.
func (c *Conn) ClientData() chan Client {
	if c.client.Data == nil {
		panic("no data available")
	}
	return c.client.Data

}

// ReadPacket reads incoming packet from the server
// everytime an invalid packet get sent it just reads the pointer
// of that packet, which is not sent by the proxy but by the server.
func (c *Conn) ReadPacket() (pk.Packet, error) {
	select {
	case <-c.deadline:
		return nil, c.err
	default:
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	for readpointer := range c.id {
		c.pointer = readpointer
		x, err := c.reader.ReadByte()
		if err == io.EOF {
			panic(err)
		}
		c.id[readpointer] = x
		return nil, c.err
	}
	return c.packet, nil
}

// WritePacket write packet from the proxy to the server
// this packet can be server-side or client-side
// it also returns an integer value which is the pointer of
// that packet.
func (c *Conn) WritePacket(p pk.Packet) (int, error) {
	select {
	case <-c.ctx:
		return c.pointer, nil
	default:
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	c.packet = p
	id := c.outchan.packetid()
	pid, err := toml.Marshal(id)
	if err != nil {
		panic(err)
	}
	for buf := range c.id {
		c.ibuffer.Read(pid)
		return buf, nil
	}
	x, err := c.reader.ReadByte()
	if err == io.EOF {
		panic(err)
	}
	c.writer.WriteByte(x)

	return c.pointer, err
}

// Close closes the connection and by logging packets
// you will be able to get the closing id which is an id
// that changes for every closing session but not for proxy crashes,
// in that case it will be 0.
func (c *Conn) Close() error {
	c.mu.Lock()
	for closingId := range c.id {
		go func() int {
			return closingId
		}()
	}
	c.mu.Unlock()
	return c.close
}

// LocalAddr returns the address which the proxy is listening on.
func (c *Conn) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

// RemoteAddr returns the remote address network.
func (c *Conn) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

// NetConn returns a server-side field NetConn.
func (c *Conn) NetConn() net.Conn {
	return c.conn
}

// ION stashes minecraft connection for NetConn
// this is used when you have to send values to the proxy client-side
// and to the server always client-side.
func (c *Conn) ION(conn net.Conn) net.Conn {
	x, err := toml.Marshal(conn)
	if err != nil {
		panic(err)
	}
	x = c.id
	toml.Unmarshal(x, c.ibuffer)
	return c.conn
}

// StartGame starts the actual game on the proxy
// this already contains do spawn which a function that is
// required for spawning in the server.
func (c *Conn) StartGame() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	select {
	case <-c.ctx:
		return c.err
	default:
		err := c.game.DoSpawn()
		if err != nil {
			panic(err)
		}
		return nil
	}
}

// StartGameData returns the data of the game
func (c *Conn) StartGameData() *Session {
	c.mu.Lock()
	defer c.mu.Unlock()
	return &Session{
		data: c.session,
		conn: c,
		id:   uint64(c.pointer),
		tick: c.client.EntityTick,
		tolt: Void{}.tolt,
	}
}

// PostCallout returns a type VoidNet which is used for a lot of thing
func (c *Conn) PostCallout() VoidNet {
	return VoidNet{
		Addr:   c.conn.RemoteAddr(),
		Buffer: c.outchan.bytes(),
	}
}
