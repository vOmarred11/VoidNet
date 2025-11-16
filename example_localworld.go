package main

import (
	"VoidNet/local"
	"VoidNet/local/level"
	"context"
)

func main() {
	cfg := local.ListenConfig{}
	listener, err := cfg.Listen(context.Background(), cfg.LanID, cfg.Signalling)
	if err != nil {
		panic(err)
	}
	for {
		err := listener.Accept()
		if err != nil {
			panic(err)
		}
		go start(listener, cfg.Conn, cfg, cfg.Player)
	}
}
func start(listener *local.Listener, conn *local.Conn, cfg local.ListenConfig, player *local.Player) {
	defer listener.Disconnect("disconnected")
	void := conn.Voidnet(listener)
	go func() {
		stash, err := void.VoidNetDecryptStash(listener.StashData())
		if err != nil {
			panic(err)
		}
		err = conn.Write(stash)
		if err != nil {
			panic(err)
		}
	}()
	conn.Wait()
	go func() {
		err := conn.StartGame(&local.Data{
			NetworkID: cfg.LanID,
			LocalData: void.OutgoingBytes,
			Game: local.Game{
				Name:       "VoidNet local",
				MaxPlayers: 10,
				MOTD:       "VoidNet LAN game",
				Spawn:      conn.NewTank().SafeSpawn(),
				Gamemode:   world.GameModeSurvival,
			},
		})
		if err != nil {
			panic(err)
		}
	}()
	conn.Wait()
	go func() {
		for {
			p, err := conn.Read(player.PlayerData())
			if err != nil {
				panic(err)
			}
			err = conn.Write(p)
			if err != nil {
				panic(err)
			}
		}
	}()
	conn.Final()
}
