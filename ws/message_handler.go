package ws

import (
	"context"
	"encoding/json"
)

// MessageHandler message handler
type MessageHandler interface {
	// HandleUnknownMessage
	HandleUnknownMessage(ctx context.Context, msg []byte)
	// HandlePushData
	HandlePushData(context.Context, PushData[ChannelRawMessage, json.RawMessage])
	// HandleResponse
	HandleResponse(context.Context, Response[ChannelRawMessage])
}

var _ MessageHandler = NoopMessageHandler{}

// NoopMessageHandler implement MessageHandler but do nonthing
type NoopMessageHandler struct{}

func (h NoopMessageHandler) HandleUnknownMessage(ctx context.Context, msg []byte) {
}

func (h NoopMessageHandler) HandlePushData(ctx context.Context, data PushData[ChannelRawMessage, json.RawMessage]) {
}

func (h NoopMessageHandler) HandleResponse(ctx context.Context, resp Response[ChannelRawMessage]) {
}
