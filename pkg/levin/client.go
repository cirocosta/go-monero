package levin

import (
	"context"
	"fmt"
	"io"
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

func (c *Client) Ping(ctx context.Context) error {
	reqHeaderB := NewRequestHeader(CommandPing, 0).Bytes()

	if _, err := c.conn.Write(reqHeaderB); err != nil {
		return fmt.Errorf("write: %w", err)
	}

	responseHeaderB := make([]byte, LevinHeaderSizeBytes)
	if _, err := io.ReadFull(c.conn, responseHeaderB); err != nil {
		return fmt.Errorf("read full header: %w", err)
	}

	fmt.Printf("resp header: %+v\n", responseHeaderB)

	// do ping
	return nil
}
