package ws

import (
	"context"

	"github.com/google/uuid"
)

// Unsubscribe
//
// https://www.okx.com/docs-v5/en/#overview-websocket-unsubscribe
func (c *Client) Unsubscribe(ctx context.Context, channels ...Channel) error {
	return write(ctx, c.getConn(), Request[Channel]{
		Id:   uuid.NewString(),
		Op:   "unsubscribe",
		Args: channels,
	})
}
