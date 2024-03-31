package ws

import (
	"context"
	"strconv"
	"time"

	"github.com/google/uuid"
	okxapigo "github.com/mind1949/okxapi-go"
)

// Login
//
// https://www.okx.com/docs-v5/en/#overview-websocket-login
func (c *Client) Login(ctx context.Context, apis ...okxapigo.Api) error {
	timestamp := time.Now().Unix()
	args := make([]loginArg, 0, len(apis))
	for _, api := range apis {
		args = append(args, loginArg{
			ApiKey:     api.Key,
			Passphrase: api.Passphrase,
			Timestamp:  strconv.FormatInt(timestamp, 10),
			Sign:       okxapigo.Sign(timestamp, api),
		})
	}
	return write(ctx, c.getConn(), Request[loginArg]{
		Id:   uuid.NewString(),
		Op:   "login",
		Args: args,
	})
}

// loginArg login arg
type loginArg struct {
	// API Key
	ApiKey string `json:"apiKey,omitempty"`
	// API Key password
	Passphrase string `json:"passphrase,omitempty"`
	// Unix Epoch time, the unit is seconds
	Timestamp string `json:"timestamp,omitempty"`
	// Signature string
	Sign string `json:"sign,omitempty"`
}

func (a loginArg) validate() error {
	// no-op
	return nil
}
