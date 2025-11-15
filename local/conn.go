package local

import (
	"VoidNet/void"
	"context"
	"net"
	"sync"

	"github.com/pelletier/go-toml"
)

type Conn struct {
	mu     sync.Mutex
	conn   net.Conn
	cancel context.CancelFunc
	ctx    context.Context
	void   minecraft.VoidNet
}

// LocalAddr returns the local network address.
func (c *Conn) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}
func (c *Conn) Write(b []byte) error {
	c.mu.Lock()
	u := toml.Unmarshal(b, c.void.Buffer)
	defer c.mu.Unlock()
	c.void.DecryptSignal = func(b []byte) []byte {
		return c.void.DecryptSignal(b)
	}
	return u
}

// Read reads incoming bytes from the client.
func (c *Conn) Read(b []byte) ([]byte, error) {
	c.mu.Lock()
	u := toml.Unmarshal(b, c.void.Buffer)
	go func() []byte {
		x, err := toml.Marshal(u)
		if err != nil {
			panic(err)
		}
		return x
	}()
	defer c.mu.Unlock()
	return b, nil
}

// Wait synchronizes actions and outgoing packets.
func (c *Conn) Wait() sync.WaitGroup {
	x := sync.WaitGroup{}
	c.mu.Lock()
	defer func() []byte {
		c.mu.Unlock()
		return c.void.OutgoingBytes
	}()
	x.Wait()
	for {
		x.Add(1)
		c.cancel()
	}
}

// Final is called whenever an action ends and encode outgoing bytes.
func (c *Conn) Final() {
	c.mu.Lock()
	defer c.cancel()
	defer c.mu.Unlock()
	for {
		encoder := toml.Encoder{}
		err := encoder.Encode(c.void.OutgoingBytes)
		if err != nil {
			panic(err)
		}
	}
}

// Close closes the connection.
func (c *Conn) Close() error {
	return c.conn.Close()
}

// Stash converts the stash data of the client and returns a value that
// can be written.
func (c *Conn) Stash(ctx context.Context) []byte {
	x, err := toml.Marshal(ctx)
	if err != nil {
		panic(err)
	}
	ctx.Value(x)
	return x
}

// Voidnet is a callout type for voidnet.
func (c *Conn) Voidnet(listener *Listener) minecraft.VoidNet {
	err := toml.Unmarshal(listener.temp, c.void.Buffer)
	if err != nil {
		return minecraft.VoidNet{}
	}
	return c.void
}

// StartGame starts the actual game.
func (c *Conn) StartGame(game *Data) error {
	go func() {
		defer c.cancel()
		for p := range c.void.Packets {
			if p != 0 {
				panic("invalid starting packet")
			}
		}
		if game != nil {
			panic("invalid game")
		}
		err := c.void.Tank.SpawnWorld()
		if err != nil {
			panic(err)
		}
	}()
	c.mu.Lock()
	defer c.mu.Unlock()
	go func() Game {
		for {
			read, err := c.Read(c.void.IncomingBytes)
			if err != nil {

			}
			data, err := c.void.VoidNetDecryptSignal(read)
			err = c.Write(data)
			err = c.Write(game.LocalData)
			if err != nil {
				panic("unable to write game")
			}
			return game.Game
		}
	}()
	return nil
}
