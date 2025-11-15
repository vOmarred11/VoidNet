package packet

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

// FeatureRegistry is a packet used to notify the client about the level generation features the server is currently
// using. This is used in combination with the client-side level generation system introduced in v1.19.20, allowing the
// client to completely generate the chunks of the level without having to rely on the server.
type FeatureRegistry struct {
	// Features is a slice of all registered level generation features.
	Features []protocol.GenerationFeature
}

// ID ...
func (pk *FeatureRegistry) ID() uint32 {
	return IDFeatureRegistry
}

func (pk *FeatureRegistry) Marshal(io protocol.IO) {
	protocol.Slice(io, &pk.Features)
}
