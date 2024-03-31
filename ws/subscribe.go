package ws

import (
	"context"

	"github.com/google/uuid"
)

// Subscribe subscribe channels
func (c *Client) Subscribe(ctx context.Context, channels ...Channel) error {
	return write(ctx, c.getConn(), Request[Channel]{
		Id:   uuid.NewString(),
		Op:   "subscribe",
		Args: channels,
	})
}
