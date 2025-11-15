package minecraft

import (
	pk "VoidNet/void/packet"
	"VoidNet/void/packs"
	"context"
	"net"

	"github.com/pelletier/go-toml"
)

// ListenConfig is the listener of the proxy
type ListenConfig struct {
	// Network returns the used network which is always TCP or VOIDNET
	// we have no clue why it changes.
	Network string
	// OnlinePlayers returns the current amount of players that is playing
	// the server.
	OnlinePlayers int32
	// MaximumPlayers returns the maximum count of players that can join
	// the server, it's still unclear why OnlinePlayers, MaximumPlayers and MOTD aren't
	// together in a status provider type.
	MaximumPlayers int32
	// RegisterDataPlayers returns the data of all the players that joined
	// at least one time in the server.
	RegisterDataPlayers []byte
	// ResourcesPacks returns the server resources packs,
	// currently on VoidNet there's no way to add custom resources packs
	// maybe it will be added in the future.
	ResourcesPacks []*packs.Pack
	// MOTD returns the motd of the server
	// if it is empty the motd will be "Minecraft Server".
	MOTD string
	// Context returns the listener context.
	Context context.Context
	// Connection returns the connection value of the proxy
	// this does not require a manual field, but it will get filled
	// only after the dialing.
	Connection *Conn
	// Call is a soft callin type for VoidNet
	// this is used for let the proxy know what values the network will have.
	Call VoidNet
}

type Listener struct {
	data       chan Client
	listener   *Listener
	context    context.Context
	raddr      string
	addr       net.Addr
	close      context.CancelFunc
	temp       []byte
	compressed []byte
}

// Listen listens on the chosen ip, this requires a network because the listener
// has to know which type should the listener be, it could only be TCP or VOIDNET
func (l *ListenConfig) Listen(network, address string) (*Listener, error) {
	l.Network = network
	nt := Network{}
	for _, n := range Networks {
		if network == n {
			nt.Listen(network, address)
		} else {
			panic("no such network")
		}
	}
	listener := &Listener{
		data:     make(chan Client),
		listener: new(Listener),
		context:  l.Context,
		raddr:    address,
	}
	return listener, nil
}

// Accept accepts the incoming request, this is the field that
// makes your screen change to "Connecting External Server" to
// "Locating External Server".
func (l *Listener) Accept() (net.Conn, error) {
	client, ok := <-l.data
	if !ok {
		return nil, &net.OpError{Op: "accept", Net: "void", Addr: l.addr, Err: net.ErrClosed}
	}
	return client.Network, nil
}

// Close closes the listener completely so it automatically closes also the connection
// and consequently disconnects also the client.
func (l *Listener) Close() error {
	return nil
}

// Disconnect disconnects the client without closing the connection
func (l *Listener) Disconnect(conn *Conn, msg string) {
	_, err := conn.WritePacket(&pk.Disconnect{
		Message: msg,
	})
	if err != nil {
		panic(err)
	}
}

// SoftCallout returns a soft callout type for the connection
// this can be used only after the StartGame action.
func (l *Listener) SoftCallout(data Listener) Conn {
	x, err := toml.Marshal(data)
	if err != nil {
		panic(err)
	}
	return Conn{
		id: x,
	}
}
