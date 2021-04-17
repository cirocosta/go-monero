package levin

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"time"
)

type Client struct {
	conn net.Conn
}

func NewClient(addr string) (*Client, error) {
	var d net.Dialer

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := d.DialContext(ctx, "tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("dial ctx: %w", err)
	}

	return &Client{
		conn: conn,
	}, nil
}

func (c *Client) Close() error {
	if c.conn == nil {
		return nil
	}

	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("close: %w", err)
	}

	return nil
}

func (c *Client) Handshake(ctx context.Context) error {
	b := NewRequest(CommandHandshake).Bytes()

	fmt.Println(b)

	if _, err := c.conn.Write(b); err != nil {
		return fmt.Errorf("write: %w", err)
	}

	resp, err := ioutil.ReadAll(c.conn)
	if err != nil {
		return fmt.Errorf("read all: %w", err)
	}

	fmt.Println(resp)
	return nil
}

func (c *Client) Ping(ctx context.Context) error {
	b := NewRequest(CommandPing).Bytes()

	if _, err := c.conn.Write(b); err != nil {
		return fmt.Errorf("write: %w", err)
	}

	resp, err := ioutil.ReadAll(c.conn)
	if err != nil {
		return fmt.Errorf("read all: %w", err)
	}

	fmt.Println(resp)

	// do ping
	return nil
}
