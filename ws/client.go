package ws

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/sync/errgroup"
)

// NewClient create okx api v5 client using websocket protocol
func NewClient(url string, fns ...OptFn) (*Client, error) {
	opts := newDefaultOPts()
	for _, fn := range fns {
		fn(&opts)
	}

	c := &Client{
		url:         url,
		opts:        opts,
		channelType: PublicChannelType,
	}
	return c, c.init()
}

// Client okx api v5 client using websocket protocol
type Client struct {
	url  string
	opts options

	channelType ChannelType
	conns       map[ChannelType]*conn

	mu sync.Mutex
	// whether or not closed
	closed atomic.Bool
}

func (c *Client) getConn() *conn {
	return c.conns[c.channelType]
}

func (c *Client) init() error {
	if len(c.conns) != 0 {
		return nil
	}
	c.conns = make(map[ChannelType]*conn)
	for _, path := range []ChannelType{
		PublicChannelType,
		PrivateChannelType,
		BusinessChannelType} {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		conn, err := newConn(ctx, c.url, path, c.opts.messageHandler)
		if err != nil {
			return err
		}
		c.conns[path] = conn
	}
	return nil
}

// Run run websocket
func (c *Client) Run(ctx context.Context) (err error) {
	eg, ctx := errgroup.WithContext(ctx)
	for _, conn := range c.conns {
		eg.Go(func() error {
			return conn.Run(ctx)
		})
	}
	return eg.Wait()
}

// Close close related resource
func (c *Client) Close() (err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.closed.CompareAndSwap(false, true) {
		return err
	}
	// close websocket connection
	for _, conn := range c.conns {
		err = errors.Join(conn.Close(), err)
	}
	return err
}

func (c *Client) Public() *Client {
	client := c.Clone()
	c.channelType = PublicChannelType
	return client
}

func (c *Client) Private() *Client {
	client := c.Clone()
	c.channelType = PrivateChannelType
	return client
}

func (c *Client) Business() *Client {
	client := c.Clone()
	c.channelType = BusinessChannelType
	return client
}

func (c *Client) Clone() *Client {
	return &Client{
		url:         c.url,
		opts:        c.opts,
		conns:       c.conns,
		channelType: c.channelType,
	}
}
