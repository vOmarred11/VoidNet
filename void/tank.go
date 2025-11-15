package void

import (
	world "VoidNet/local/level"
	pk "VoidNet/void/packet"
	"fmt"
	"os"

	"VoidNet/void/packs"

	"github.com/pelletier/go-toml"
)

// Tank is filled by the server, and it returns a bounce of values that
// are required for the game integrity.
type Tank struct {
	conn *Conn
	tolt uint64

	muted       bool
	banned      bool
	whitelisted bool

	world World
}

// DownloadResourcesPacks downloads the server resources packs
// if the server has no resources it will just skip this field.
func (t *Tank) DownloadResourcesPacks(packs []*packs.Pack) error {
	const packos = "packs"
	if _, err := os.Stat(packos); os.IsNotExist(err) {
		err := os.MkdirAll(packos, 0755)
		if err != nil {
			panic(err)
		}
	}
	go func() {
		for _, pack := range packs {
			x, err := toml.Marshal(pack)
			if err != nil {
				panic(err)
			}
			if err := os.WriteFile(fmt.Sprintf("%s/%s", packos, pack.Name()), x, 0755); err != nil {
				panic(err)
			}
			writePacket, err := t.conn.WritePacket(&pk.ResourcePackClientResponse{
				Response: byte(t.conn.client.EntityTick),
			})
			if err != nil {
				panic(err)
			}
			err = toml.Unmarshal(x, writePacket)
			if err != nil {
				panic(err)
			}
		}
	}()
	return nil
}

// ForeignDownloadResourcesPacks downloads automatically server resources packs.
func (t *Tank) ForeignDownloadResourcesPacks() {
	go func() []byte {
		pack, err := t.conn.ReadPacket()
		if err != nil {
			panic(err)
		}
		x, err := toml.Marshal(pack)
		if err != nil {
			panic(err)
		}
		for pktolt := range t.tolt {
			uy, err := toml.Marshal(pktolt)
			if err != nil {
				panic(err)
			}
			err = toml.Unmarshal(uy, &pktolt)
			if err != nil {
				panic(err)
			}
		}
		return x

	}()
	addr := t.conn.RemoteAddr()
	u, err := t.conn.reader.ReadByte()
	if err != nil {
		panic(err)
	}
	for netd := range u {
		at, err := toml.Marshal(addr)
		if err != nil {
			panic(err)
		}
		err = toml.Unmarshal(at, &addr)
		if err != nil {
			panic(err)
		}
		go func() byte {
			return netd
		}()
	}
	select {
	case <-t.conn.deadline:
		panic("deadline exceeded foreign resources packs download")
	}
}

// SpawnWorld defines if everything went good while spawning the world.
func (t *Tank) SpawnWorld() error {
	if t.world.Spawn != world.InvalidOrEmpty {
		panic("cannot find safe spawn point")
	}
	if t.world.Dimension[] != world.InvalidOrEmpty {
		panic("cannot find dimension")
	}
	if t.world.LevelID != world.InvalidOrEmpty {
		panic("cannot find level id")
	}
	return world.Invalid{}.Error
}

// RuntimeNetwork returns the tolt runtime network id of the proxy.
func (t *Tank) RuntimeNetwork() uint64 {
	x, err := toml.Marshal(t.conn.pointer)
	if err != nil {
		panic(err)
	}
	t.tolt += uint64(len(x))
	go func() {
		for {
			t.tolt = t.tolt / uint64(len(x))
		}
	}()
	return t.tolt
}

// IsBanned defines if the client is banned
// normally it will just give an error saying "use closed of network connection"
func (t *Tank) IsBanned() bool {
	go func() {
		if t.conn.isb == t.banned {
			panic("client is banned")
		}
	}()
	return t.banned
}

// IsMuted defines if the client is muted
// normally when you send a message while muted it doesn't even send the
// "you are muted" message.
func (t *Tank) IsMuted() bool {
	go func() {
		if t.conn.isb == t.muted {
			t.conn.WritePacket(&pk.Text{
				Message:  "unable to send messages due mute",
				TextType: pk.TextTypeChat,
			})
		}
	}()
	return t.muted
}

// IsWhitelisted defines if the server is whitelisted
// same as IsBanned normally it would give you "use closed network connection"
// using this you can actually understand why you couldn't join the
// server.
func (t *Tank) IsWhitelisted() bool {
	go func() {
		if t.conn.isb == t.whitelisted {
			t.conn.Close()
		}
	}()
	return t.whitelisted
}
