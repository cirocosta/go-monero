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

func (c *Client) Handshake(ctx context.Context) error {
	payload := (&PortableStorage{
		Entries: []Entry{
			{
				Name: "node_data",
				Serializable: &Section{
					Entries: []Entry{
						{
							Name:         "local_time",
							Serializable: BoostUint64(time.Now().Unix()),
						},
						{
							Name:         "my_port",
							Serializable: BoostUint32(0),
						},
						{
							Name:         "network_id",
							Serializable: BoostString(string(MainnetNetworkId)),
						},
						{
							Name:         "peer_id",
							Serializable: BoostUint64(12343332),
						},
					},
				},
			},

			{
				Name: "payload_data",
				Serializable: &Section{
					Entries: []Entry{
						{
							Name:         "cumulative_difficulty",
							Serializable: BoostUint64(1),
						},
						{
							Name:         "current_height",
							Serializable: BoostUint64(1),
						},
						{
							Name:         "top_id",
							Serializable: BoostString(MainnetGenesisTx),
						},
						{
							Name:         "top_version",
							Serializable: BoostByte(1),
						},
					},
				},
			},
		},
	}).Bytes()

	reqHeaderB := NewRequestHeader(CommandHandshake, uint64(len(payload))).Bytes()

	if _, err := c.conn.Write(reqHeaderB); err != nil {
		return fmt.Errorf("write header: %w", err)
	}

	if _, err := c.conn.Write(payload); err != nil {
		return fmt.Errorf("write payload: %w", err)
	}

	responseHeaderB := make([]byte, LevinHeaderSizeBytes)
	if _, err := io.ReadFull(c.conn, responseHeaderB); err != nil {
		return fmt.Errorf("read full header: %w", err)
	}

	respHeader, err := NewHeaderFromResponseBytes(responseHeaderB)
	if err != nil {
		return fmt.Errorf("new header from resp bytes: %w", err)
	}

	fmt.Printf("%+v\n", respHeader)

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

	respHeader, err := NewHeaderFromResponseBytes(responseHeaderB)
	if err != nil {
		return fmt.Errorf("new header from resp bytes: %w", err)
	}

	fmt.Printf("%+v\n", respHeader)

	return nil
}
