package main

import (
	"VoidNet/void"
)

func main() {
	cfg := void.ListenConfig{}
	listener, err := cfg.Listen("voidnet", "0.0.0.0:19132")
	if err != nil {
		panic(err)
	}
	c, err := listener.Accept()
	if err != nil {
		panic(err)
	}
	dialer, err := (&void.VoidDialer{
		Connection:   cfg.Connection,
		CallAccept:   c,
		ClientData:   cfg.Connection.ClientData(),
		Callin:       cfg,
		Callout:      cfg.Call,
		VoidCallType: void.ByteDataTrian,
	}).Dial(&cfg.Call, "play.lbsg.net:19132")
	tank := dialer.PostCallout().Tank
	err = tank.DownloadResourcesPacks(cfg.ResourcesPacks)
	if err != nil {
		panic(err)
	}
	cfg.Call.ByteCall = func(data uint8) uint8 {
		return void.ByteData
	}
	go func() {
		err := dialer.StartGame()
		if err != nil {
			panic(err)
		}
		err = tank.SpawnWorld()
		if err != nil {
			panic(err)
		}
	}()
	go func() {
		defer dialer.Close()
		defer listener.Disconnect(dialer, "closed")
		for {
			packet, err := dialer.ReadPacket()
			if err != nil {
				panic(err)
			}
			if err, _ := dialer.WritePacket(packet); err != 0 {
				return
			}
		}
	}()
}
