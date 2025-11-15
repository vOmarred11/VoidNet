package local

import (
	"archive/zip"
	"os"
	"sync"

	"github.com/pelletier/go-toml"
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
		t.conn.Write(x)
		t.mu.Lock()
	}()
	defer t.mu.Unlock()
	return t.conn.void.OutgoingBytes
}
