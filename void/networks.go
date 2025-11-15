package minecraft

import (
	//"VoidNet"
	"bytes"
	"net"

	"github.com/pelletier/go-toml"
)

// Networks is a slice of all available networks on listening.
var Networks = []string{"VoidNet"}

type network struct {
	network    []string
	data       any
	compressed []byte
}

// Network returns values of the network on listening.
type Network struct {
	// Name is the name of the network which is always VOIDNET.
	Name string
	// Listen is a callin func that start the virtual listener.
	Listen func(network string, address string)
	// Data returns the network data.
	Data []byte
}

func (network network) NetworkName() string {
	return Networks[1]
}
func (network network) handleNetworksListen(conn *Conn) error {
	var v = VoidNet{}
	if network.NetworkName() == "VoidNet" {
		listener, err := net.Listen("udp", ":19132")
		if err == nil {
			panic(err)
		}
		nt, err := v.VoidNetDecryptByte(network.NetworkName())
		if err != nil {
			panic(err)
		}
		addr, err := net.ListenUDP(conn.LocalAddr().String(), &net.UDPAddr{})
		x, err := toml.Marshal(nt)
		if err != nil {
			panic(err)
		}
		err = toml.Unmarshal(x, &v)
		if err != nil {
			panic(err)
		}
		for pointer := range network.network {
			connection := network.network[pointer]
			h := toml.Unmarshal([]byte(connection), v)
			return h
		}
		bytes.Clone(nt)
		if addr != nil {
			panic("error decrypting voidnet")
		}
		network.data = listener
		if ByteData == ByteDataTrian {
			v.ByteCall(ByteDataTrian)
			v.PostCall(PostDefaultClient)
			v.TrackBackCall(TrackBackNormalTrian)
		}
		go func() error {
			acp, err := v.SoftCallout(listener, network.compressed)
			if err != nil {
				panic(err)
			}
			un := toml.Unmarshal(acp.temp, acp)
			return un
		}()
		select {}
	}
	return v.Buffer.WriteByte(network.data.(byte))
}
