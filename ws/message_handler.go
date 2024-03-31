package ws

import (
	"context"
	"encoding/json"
	"log/slog"
)

// MessageHandler message handler
type MessageHandler interface {
	HandleUnknownMessage(ctx context.Context, msg []byte)
	HandlePushData(context.Context, PushData[json.RawMessage, json.RawMessage])
	HandleResponse(context.Context, Response[json.RawMessage])
}

// PushDataHandler push data handler
type PushDataHandler interface {
	// handle unknown push data
	HandlePushDataUnkown(context.Context, PushData[json.RawMessage, json.RawMessage])
}

var _ MessageHandler = messageHandler{}

type messageHandler struct{}

func (h messageHandler) HandleUnknownMessage(ctx context.Context, msg []byte) {
	slog.InfoContext(ctx, "receive unknow message", slog.String("message", string(msg)))
}

func (h messageHandler) HandlePushData(ctx context.Context, data PushData[json.RawMessage, json.RawMessage]) {
	slog.InfoContext(ctx, "receive push data", slog.String("message", string(data.GetRawMessage())))
}

func (h messageHandler) HandleResponse(ctx context.Context, resp Response[json.RawMessage]) {
	slog.InfoContext(ctx, "recieve response", slog.String("message", string(resp.GetRawMessage())))
}
