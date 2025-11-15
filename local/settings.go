package local

import (
	world "VoidNet/local/level"

	"github.com/pelletier/go-toml"
)

// GameSettings are all the settings of the client that influxes the game.
type GameSettings struct {
	// player returns the client data.
	player *Conn

	// Gamemode is the current client gamemode.
	Gamemode world.GameMode
	// Difficulty is the current world difficulty.
	Difficulty world.Difficulty
	// Hardcore defines if the world is in the hardcore mode.
	Hardcore bool

	// Seed returns the seed on the world
	// this one will be empty if there will be some errors in the
	// world generation.
	Seed uint32
	// FlatWorld defines if the world is flat.
	FlatWorld bool
	// ShowCoordinates shows current player position at the top.
	ShowCoordinates bool
	// ShowDaysPlayed defines how much minecraft days you've played since
	// the start of the world
	// for some reason if DayLightCycle is set to false it will continue
	// counting the days.
	ShowDaysPlayed bool
	// RecipeUnlocking defines if when collecting materials it should unlock recipes.
	RecipeUnlocking bool
	// FireSpreads defines if fire can spread between entities.
	FireSpreads bool
	// TNTExplodes defines if tnt should explode.
	TNTExplodes bool
	// MobLoot defines if mob should drop loot when they die.
	MobLoot bool
	// NaturalGeneration defines if health should regen itself.
	NaturalGeneration bool
	// TileDrops defines if it should drop blocks when they are broken.
	TileDrops bool
	// SkipNightBySleeping defines if when you sleep it should skip the night.
	SkipNightBySleeping bool
	// RequiredSleepingPlayers is the amount of players that are required to
	// skip the night by sleeping.
	RequiredSleepingPlayers int
	// ImmediateRespawn defines if when you die the client will respawn instantly.
	ImmediateRespawn bool
	// RespawnBlocksExplode defines if blocks like Respawn Anchor should explode
	// if they are not in their dimension.
	RespawnBlocksExplode bool
	// RespawnRadius is the radius of your respawn point (Max: 128).
	RespawnRadius int32

	// Cheats defines is cheats are enabled in this world.
	Cheats bool
	// DayLightCycle defines if the normal light cycle should continue.
	DayLightCycle bool
	// Time defines the current time right now.
	Time int
	// KeepInventory defines that it should keep you inventory when you die.
	KeepInventory bool
	// MobSpawning defines if mob should spawn.
	MobSpawning bool
	// MobGrieffing defines if mob should destroy blocks or else on you world
	// a clear example can be the Creeper.
	MobGrieffing bool
	// EntitiesDropLoot defines if entities such a Cow or a Sheep when they died they should drop loot.
	EntitiesDropLoot bool
	// WhetherCycle defines if whether should change.
	WhetherCycle bool

	// BetaAPIs defines if beta test versions are enabled
	// I recommend keep this off if you are on a listener.
	BetaAPIs bool
}

// ApplySettings applies current changes.
func (g GameSettings) ApplySettings() error {
	x, err := toml.Marshal(g)
	if err != nil {
		panic(err)
	}
	return g.player.Write(x)
}

// DefaultSettings returns the default minecraft settings of the world.
func (g GameSettings) DefaultSettings() GameSettings {
	game := GameSettings{
		player:                  g.player,
		Gamemode:                world.GameModeSurvival,
		Difficulty:              world.DifficultyNormal,
		Hardcore:                false,
		Seed:                    g.Seed,
		FlatWorld:               false,
		ShowCoordinates:         true,
		ShowDaysPlayed:          false,
		RecipeUnlocking:         true,
		FireSpreads:             true,
		TNTExplodes:             true,
		MobLoot:                 true,
		NaturalGeneration:       true,
		TileDrops:               true,
		SkipNightBySleeping:     true,
		RequiredSleepingPlayers: 0,
		ImmediateRespawn:        false,
		RespawnBlocksExplode:    false,
		RespawnRadius:           5,
		Cheats:                  false,
		DayLightCycle:           true,
		Time:                    world.TimeDay,
		KeepInventory:           false,
		MobSpawning:             true,
		MobGrieffing:            true,
		EntitiesDropLoot:        true,
		WhetherCycle:            true,
		BetaAPIs:                false,
	}
	return game
}
