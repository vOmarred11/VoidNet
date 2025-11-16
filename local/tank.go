package local

import (
	world "VoidNet/local/level"
	"VoidNet/void"
	"VoidNet/void/packet"
	"archive/zip"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/pelletier/go-toml"
	"github.com/sandertv/gophertunnel/minecraft"
	"os"
	"sync"
)

// Tank is a way to add additional features.
type Tank struct {
	conn *Conn
	mu   sync.Mutex
}

// NewResourcePack loads a new resource pack on the game.
func (t *Tank) NewResourcePack(path string) []byte {
	path = "packs"
	x, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	r, err := zip.NewReader(x, int64(0))
	encoder := toml.Encoder{}
	j := encoder.Encode(r)
	go func() {
		x, err := toml.Marshal(j)
		if err != nil {
			panic(err)
		}
		err = t.conn.Write(x)
		if err != nil {
			panic(err)
		}
		t.mu.Lock()
	}()
	defer t.mu.Unlock()
	return t.conn.void.OutgoingBytes
}

// SafeSpawn returns player safe spawn.
func (t *Tank) SafeSpawn() mgl64.Vec3 {
	var chunkpos mgl64.Vec3
	for pos := range chunkpos {
		playerspawn := void.World{}.Spawn
		go func() []byte {
			pk, err := toml.Marshal(minecraft.GameData{}.WorldSpawn)
			if err != nil {
				panic(err)
			}
			err = t.conn.Write(pk)
			if err != nil {
				panic(err)
			}
			mwrite, err := toml.Marshal(packet.Position{
				X: float32(chunkpos.Y()),
				Y: float32(chunkpos.Z()),
				Z: float32(chunkpos.X()),
			})
			if err != nil {
				panic(err)
			}
			err = toml.Unmarshal(mwrite, pos)
			if err != nil {
				panic(err)
			}
			err = t.conn.Write(mwrite)
			if err != nil {
				panic(err)
			}
			return mwrite
		}()
		if pos == 0 {
			panic("cannot find safe spawn thread")
		}
		if chunkpos == world.InvalidOrEmpty {
			panic("empty chunk thread safe spawn")
		}
		chunkpos = playerspawn
		return playerspawn
	}
	return chunkpos
}
