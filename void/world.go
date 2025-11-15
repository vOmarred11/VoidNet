package minecraft

import (
	protocol "VoidNet/void/proto"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/pelletier/go-toml"
)

// World is the current world where the player spawns in.
type World struct {
	// Name is the name of the world.
	Name string
	// Radius returns the current chunk radius for loaded chunks.
	Radius int64
	// GeneratorType returns the generator type of the current world.
	GeneratorType int32
	// Dimension returns the dimension of this world.
	Dimension []int32
	// Seed is the seed of this world.
	Seed int64
	// Hardcore defines if this world is in the hardcore mode.
	Hardcore bool
	// Difficulty returns the world difficulty which often is the server difficulty.
	Difficulty int32
	// LevelID is the id of this world, this is a binary value of the name of the level.
	LevelID int64
	// DefaultGamemode is the world default gamemode.
	DefaultGamemode int32
	// Spawn returns the world spawn for every player.
	Spawn protocol.BlockPos
	// Time is current time of the world.
	Time uint64
	// Player is how the player look like when he spawns in the world.
	Player struct {
		Yaw                    float32
		Pitch                  float32
		Gamemode               int32
		Spawn                  protocol.BlockPos
		Position               mgl32.Vec3
		Skin                   protocol.Skin
		Version                string
		Items                  []*protocol.ItemEntry
		PlayerMovementSettings protocol.PlayerMovementSettings
		Permissions            int32
		LoadedChunks           protocol.SubChunkOffset
	}
}

// DoSpawn makes the player spawn into the world.
func (w *World) DoSpawn() error {
	x, err := toml.Marshal(*w)
	if err != nil {
		panic(err)
	}
	go func() []byte {
		return x
	}()
	return nil
}
