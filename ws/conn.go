package ws

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"sync/atomic"
	"time"

	"golang.org/x/sync/errgroup"
	"nhooyr.io/websocket"
)

// ChannelType
type ChannelType string

const (
	PublicChannelType   = "public"
	PrivateChannelType  = "private"
	BusinessChannelType = "business"
)

// newConn new conn
func newConn(ctx context.Context, baseURL string, path ChannelType, handler MessageHandler) (*conn, error) {
	if handler == nil {
		handler = NoopMessageHandler{}
	}
	c := &conn{
		baseURL: baseURL,
		path:    string(path),
		handler: handler,
	}
	return c, c.init(ctx)
}

// conn connection
//
// https://www.okx.com/docs-v5/en/#overview-websocket-connect
type conn struct {
	baseURL string
	path    string
	url     string

	// used websocket connection
	*websocket.Conn
	// last activity's unix time (the number of milliseconds)
	lastActive atomic.Int64

	handler MessageHandler

	reqId      atomic.Uint64
	respId     atomic.Uint64
	pushDataId atomic.Uint64

	closed atomic.Bool
}

// Run run conn
func (c *conn) Run(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return c.periodicPing(ctx)
	})
	eg.Go(func() error {
		return c.consume(ctx)
	})
	return eg.Wait()
}

func (c *conn) init(ctx context.Context) error {
	u, err := url.Parse(c.baseURL)
	if err != nil {
		return err
	}
	u = u.JoinPath(c.path)
	c.url = u.String()
	c.Conn, err = dial(ctx, u.String())
	if err != nil {
		return err
	}
	return nil
}

func (c *conn) Read(ctx context.Context) (websocket.MessageType, []byte, error) {
	typ, b, err := c.Conn.Read(ctx)
	if err != nil {
		return typ, b, err
	}
	c.refreshLastActive()
	return typ, b, nil
}

func (c *conn) Write(ctx context.Context, typ websocket.MessageType, p []byte) error {
	err := c.Conn.Write(ctx, typ, p)
	if err != nil {
		return err
	}
	c.refreshLastActive()
	return nil
}

// periodicPing periodic ping
//
// https://www.okx.com/docs-v5/en/#overview-websocket-connect
func (c *conn) periodicPing(ctx context.Context) error {
	timeout := (30 - 5) * time.Second
	ticker := time.NewTicker(timeout)
	defer ticker.Stop()
	for i := 0; ; i++ {
		select {
		case <-ctx.Done():
			return context.Cause(ctx)
		case <-ticker.C:
			// no-op
		}

		if c.closed.Load() {
			return nil
		}

		err := c.Write(ctx, websocket.MessageText, []byte("ping"))
		if err != nil {
			slog.ErrorContext(ctx, "ping failed",
				slog.Group("okxapi",
					slog.String("err", err.Error()),
					slog.String("url", c.url),
				))
			continue
		}
		slog.InfoContext(ctx, "ping successful",
			slog.Group("okxapi",
				slog.String("url", c.url),
			))
	}
}

func (c *conn) consume(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return context.Cause(ctx)
		default:
			// no-op
		}

		msg, rawMsg, err := read[message](ctx, c)
		if err != nil {
			slog.ErrorContext(ctx, "read websocket message failed",
				slog.Any("error", err))
			continue
		}
		if string(rawMsg) == "pong" {
			continue
		}
		switch msg.Event {
		case "":
			c.handler.HandlePushData(ctx, PushData[ChannelRawMessage, json.RawMessage]{
				Arg:        msg.Arg,
				Data:       msg.Data,
				rawMessage: rawMsg,
			})
		case "subscribe", "unsubscribe", "error":
			c.handler.HandleResponse(ctx, message2Response(*msg))
		default:
			c.handler.HandleUnknownMessage(ctx, rawMsg)
		}
	}
}

// message2Response
func message2Response(msg message) Response[ChannelRawMessage] {
	return Response[ChannelRawMessage]{
		Event:  msg.Event,
		Arg:    msg.Arg,
		Code:   msg.Code,
		Msg:    msg.Msg,
		ConnId: msg.ConnId,
	}
}

// refreshLastActive refresh conn's last activity time
func (c *conn) refreshLastActive() {
	c.lastActive.Swap(time.Now().UnixMilli())
}

// Close close conn
func (c *conn) Close() (err error) {
	if c == nil || c.Conn == nil {
		return nil
	}
	c.closed.Swap(true)
	return c.Conn.Close(websocket.StatusNormalClosure, "")
}

// dial
func dial(ctx context.Context, url string) (*websocket.Conn, error) {
	conn, _, err := websocket.Dial(ctx, url, nil)
	if err != nil {
		const format = "dial websocket server(addr: %q) failed"
		return nil, errors.Join(fmt.Errorf(format, url), err)

	}
	return conn, err
}

// write
func write[T any](ctx context.Context, conn *conn, req Request[T]) (err error) {
	defer func() {
		if err != nil {
			err = errors.Join(fmt.Errorf("failed to write"))
		}
	}()
	if conn == nil {
		return fmt.Errorf("websocket connection is nil")
	}

	err = validate(req.Args...)
	if err != nil {
		return err
	}
	b, err := json.Marshal(&req)
	if err != nil {
		return err
	}
	return conn.Write(ctx, websocket.MessageText, b)
}

type message struct {
	// Operation
	// login
	// error
	Event string `json:"event"`

	Arg  ChannelRawMessage `json:"arg"`
	Data []json.RawMessage `json:"data"`
	// Error code
	Code string `json:"code"`
	// Error message
	Msg string `json:"msg"`
	// WebSocket connection ID
	ConnId string `json:"connId"`
}

// read read from conn
func read[Message any](ctx context.Context, conn *conn) (_ *Message, _ []byte, err error) {
	for {
		select {
		case <-ctx.Done():
			return nil, nil, context.Cause(ctx)
		default:
			// no-op
		}

		var msg Message
		typ, b, err := conn.Read(ctx)
		if err != nil {
			return nil, nil, err
		}
		if typ != websocket.MessageText {
			slog.InfoContext(ctx, fmt.Sprintf("got %s type message", typ))
			continue
		}
		if string(b) == "pong" {
			return nil, b, nil
		}
		err = json.NewDecoder(bytes.NewReader(b)).Decode(&msg)
		if err != nil {
			return nil, nil, err
		}
		return &msg, b, nil
	}
}
