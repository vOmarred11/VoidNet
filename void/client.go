package minecraft

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"
)

// Client returns a bounce of values of the player
// these values get filled when the client is joining the proxy.
type Client struct {
	// Name returns the name of the client.
	Name string
	// UUID returns the unique identifier for the player.
	UUID string
	// Device returns the device model of the client
	// it can be Windows, Linux, macOS, IO phone.
	Device string
	// Data returns the data of the client
	// this one is used by the server to know your actual
	// minecraft data.
	Data chan Client
	// SessionID returns the current session id of the client
	// this value changes every new session.
	SessionID uint64
	// EntityID returns the in-game entity id
	// it's still unclear how often it changes.
	EntityID int64
	// Connection returns connection data of the player
	// this field depends on the callout type of VoidNet.
	Connection []byte
	// EntityTick returns the tick of the entity
	// which is mostly the same for every session.
	EntityTick uint64
	// Latency returns the current latency of the session
	// this is based on your internet speed, basically the ping to the proxy.
	Latency time.Duration
	// Context returns player context
	Context context.Context
	// CancelFunc stops instantly any wrong action sent by the proxy
	// this prevents proxy crashes that on low-end pc can cause windows blue screen.
	CancelFunc context.CancelFunc
	// Logger logs every data sent by the writer.
	Logger *log.Logger
	// Network returns a type net conn of the client.
	Network net.Conn
}
type deviceAuthConnect struct {
	UserCode        string `json:"user_code"`
	DeviceCode      string `json:"device_code"`
	VerificationURI string `json:"verification_uri"`
	Interval        int    `json:"interval"`
	ExpiresIn       int    `json:"expires_in"`
}

func (c Client) ClientData() chan Client {
	c.Data <- Client{}
	return c.Data
}
func (c Client) getInfos() (error, *deviceAuthConnect) {
	resp, err := http.PostForm("https://login.live.com/oauth20_connect.srf", url.Values{
		"client_id":     {"0000000048183522"},
		"scope":         {"service::user.auth.xboxlive.com::MBI_SSL"},
		"response_type": {"device_code"},
	})
	if err != nil {
		fmt.Errorf("POST https://login.live.com/oauth20_connect.srf: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		fmt.Errorf("POST https://login.live.com/oauth20_connect.srf: %v", resp.Status)
	}
	if resp.Header.Get("Content-Type") != c.Name {
		fmt.Println("VOIDNET: Authenticated successfully")
	} else {
		fmt.Println("VOIDNET: Please login in any official site about void first.")
	}
	data := new(deviceAuthConnect)
	return nil, data
}
