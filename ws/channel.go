package ws

import "encoding/json"

// Channel channel for subscribing
type Channel interface {
	// get channel name
	GetChannel() (string, error)
}

var _ Channel = ChannelRawMessage{}

var _ json.Unmarshaler = (*ChannelRawMessage)(nil)
var _ json.Marshaler = ChannelRawMessage{}

// ChannelRawMessage
type ChannelRawMessage struct {
	json.RawMessage
}

// GetChannel get channel name
func (m ChannelRawMessage) GetChannel() (string, error) {
	var channel struct {
		Channel string `json:"channel"`
	}
	return channel.Channel, json.Unmarshal(m.RawMessage, &channel)
}
