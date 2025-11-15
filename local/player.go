package local

import "github.com/pelletier/go-toml"

// Player is the entity of the world
type Player struct {
	// PlayerID is the local identifier of the player.
	PlayerID uint64
	// Tolt is the tick that responses from the client to the host.
	Tolt     uint8
}

// PlayerData returns the data of the player in a byte type.
func (e Player) PlayerData() []byte {
	x, err := toml.Marshal(e)
	if err != nil {
		panic(err)
	}
	return x
}
