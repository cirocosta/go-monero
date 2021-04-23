package levin

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"time"
)

const DialTimeout = 15 * time.Second

type Client struct {
	conn net.Conn
}

func NewClient(addr string) (*Client, error) {
	var d net.Dialer

	ctx, cancel := context.WithTimeout(context.Background(), DialTimeout)
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

func (c *Client) Handshake(ctx context.Context) (*LocalPeerList, error) {
	payload := (&PortableStorage{
		Entries: []Entry{
			{
				Name: "node_data",
				Serializable: &Section{
					Entries: []Entry{
						{
							Name:         "network_id",
							Serializable: BoostString(string(MainnetNetworkId)),
						},
					},
				},
			},
		},
	}).Bytes()

	reqHeaderB := NewRequestHeader(CommandHandshake, uint64(len(payload))).Bytes()

	if _, err := c.conn.Write(reqHeaderB); err != nil {
		return nil, fmt.Errorf("write header: %w", err)
	}

	if _, err := c.conn.Write(payload); err != nil {
		return nil, fmt.Errorf("write payload: %w", err)
	}

again:
	responseHeaderB := make([]byte, LevinHeaderSizeBytes)
	if _, err := io.ReadFull(c.conn, responseHeaderB); err != nil {
		return nil, fmt.Errorf("read full header: %w", err)
	}

	respHeader, err := NewHeaderFromBytesBytes(responseHeaderB)
	if err != nil {
		return nil, fmt.Errorf("new header from resp bytes: %w", err)
	}

	dest := new(bytes.Buffer)

	if respHeader.Length != 0 {
		if _, err := io.CopyN(dest, c.conn, int64(respHeader.Length)); err != nil {
			return nil, fmt.Errorf("copy payload to stdout: %w", err)
		}
	}

	if respHeader.Command != CommandHandshake {
		dest.Reset()
		goto again
	}

	ps, err := NewPortableStorageFromBytes(dest.Bytes())
	if err != nil {
		return nil, fmt.Errorf("new portable storage from bytes: %w", err)
	}

	peerList := NewLocalPeerListFromEntries(ps.Entries)
	return &peerList, nil
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

	respHeader, err := NewHeaderFromBytesBytes(responseHeaderB)
	if err != nil {
		return fmt.Errorf("new header from resp bytes: %w", err)
	}

	fmt.Printf("%+v\n", respHeader)

	return nil
}
