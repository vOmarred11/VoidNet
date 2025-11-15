package local

import (
	"VoidNet/local/level"
	"VoidNet/void/packs"
	protocol "VoidNet/void/proto"

	"github.com/pelletier/go-toml"
	"github.com/sandertv/gophertunnel/minecraft"
)

// Game is a bounce of values when the game starts.
type Game struct {
	// Name is the name of you local world
	// if this is not set it will be "default".
	Name string
	// Gamemode is the default gamemode of the world
	// if this is not set the gamemode will be survival.
	Gamemode world.GameMode
	// Spawn it the spawn of the player in the world
	// if this not set it will be {0 , 0 , 0}
	// I recommend setting this otherwise you will spawn in the void.
	Spawn float64
	// MaxPlayers is the max count players that can join this world.
	MaxPlayers int
	// Players is the current amount of players.
	Players int
	// MOTD is the status provider of the local world
	// if this is not set it will be "Minecraft Server".
	MOTD string
	// ResourcesPacks is a slice of all resources packs in the world.
	ResourcesPacks []*packs.Pack
	// Skin returns the skin data of the player.
	Skin protocol.Skin
	// Protocols is a slice of the minecraft protocol.
	Protocols []minecraft.Protocol
	// Items is a slice of all the item in the world,
	// custom won't be found it this slice.
	Items []*protocol.ItemStack
}

// Data is the data of the local network.
type Data struct {
	// NetworkID returns the network id.
	NetworkID uint64
	// LocalData returns the connection dat.
	LocalData []byte
	// Game returns game settings.
	Game Game
}

// StartGame starts the actual game on the local.
func (d *Data) StartGame(prop Game) error {
	d.defaultValues(prop)
	x := world.New()
	prop.Name = x.Name()
	prop.Spawn = float64(x.Spawn().X() + x.Spawn().Y())
	prop.Gamemode = x.DefaultGameMode()
	go func() {
		defer func() []byte {
			x, err := toml.Marshal(d)
			if err != nil {
				panic(err)
			}
			return x
		}()
	}()
	x.Save()
	return nil
}

// defaultValues load default values.
func (d *Data) defaultValues(prop Game) {
	if prop.Name == "" {
		prop.Name = "default"
	}
	if prop.MaxPlayers == 0 {
		prop.MaxPlayers = 10
	}
	if prop.MOTD == "" {
		prop.MOTD = "LAN game"
	}
	if prop.Gamemode == world.InvalidOrEmpty {
		prop.Gamemode = world.GameModeSurvival
	}
	if prop.Items == nil {
		prop.Items = []*protocol.ItemStack{}
	}
}
