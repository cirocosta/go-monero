package zmq

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-zeromq/zmq4"
)

type Client struct {
	endpoint string
	topic    Topic
	sub      zmq4.Socket
}

// NewClient instantiates a new client that will receive monerod's zmq events.
//
// 	- `topic` is a fully-formed zmq topic to subscribe to
//
// 	- `endpoint` is the full address where monerod has been configured to
// 	publish the messages to, including the network schama. for instance,
// 	considering that monerod has been started with
//
//		monerod --zmq-pub tcp://127.0.0.1:18085
//
//	`endpoint` should be 'tcp://127.0.0.1:18085'.
//
//
func NewClient(endpoint string, topic Topic) *Client {
	return &Client{
		endpoint: endpoint,
		topic:    topic,
	}
}

// Stream provides channels where instances of the desired topic object are
// sent to.
//
type Stream struct {
	ErrC chan error

	FullChainMainC    chan *FullChainMain
	FullTxPoolAddC    chan *FullTxPoolAdd
	MinimalChainMainC chan *MinimalChainMain
	MinimalTxPoolAddC chan *MinimalTxPoolAdd
}

// Listen listens for a topic pre-configured for this client (via NewClient).
//
// Clients listen to a single topic at a time - to listen to multiple topics,
// create new clients and listen on the corresponding stream's channel.
//
func (c *Client) Listen(ctx context.Context) (*Stream, error) {
	if err := c.listen(ctx, c.topic); err != nil {
		return nil, fmt.Errorf("listen on '%s': %w", c.topic, err)
	}

	stream := &Stream{
		ErrC: make(chan error),

		FullChainMainC:    make(chan *FullChainMain),
		FullTxPoolAddC:    make(chan *FullTxPoolAdd),
		MinimalChainMainC: make(chan *MinimalChainMain),
		MinimalTxPoolAddC: make(chan *MinimalTxPoolAdd),
	}

	go func() {
		if err := c.loop(stream); err != nil {
			stream.ErrC <- fmt.Errorf("loop: %w", err)
		}

		close(stream.ErrC)

		close(stream.FullChainMainC)
		close(stream.FullTxPoolAddC)
		close(stream.MinimalChainMainC)
		close(stream.MinimalTxPoolAddC)
	}()

	return stream, nil
}

// Close closes any established connection, if any.
//
func (c *Client) Close() error {
	if c.sub == nil {
		return nil
	}

	return c.sub.Close()
}

func (c *Client) listen(ctx context.Context, topic Topic) error {
	c.sub = zmq4.NewSub(ctx)

	err := c.sub.Dial(c.endpoint)
	if err != nil {
		return fmt.Errorf("dial '%s': %w", c.endpoint, err)
	}

	err = c.sub.SetOption(zmq4.OptionSubscribe, string(topic))
	if err != nil {
		return fmt.Errorf("subscribe: %w", err)
	}

	return nil
}

func (c *Client) loop(stream *Stream) error {
	for {
		msg, err := c.sub.Recv()
		if err != nil {
			return fmt.Errorf("recv: %w", err)
		}

		for _, frame := range msg.Frames {
			err := c.ingestFrameArray(stream, frame)
			if err != nil {
				return fmt.Errorf("consume frame: %w", err)
			}
		}
	}
}

func (c *Client) ingestFrameArray(stream *Stream, frame []byte) error {
	topic, gson, err := jsonFromFrame(frame)
	if err != nil {
		return fmt.Errorf("json from frame: %w", err)
	}

	if c.topic != topic {
		return fmt.Errorf("topic '%s' doesn't match "+
			"expected '%s'", topic, c.topic)
	}

	switch c.topic {
	case TopicFullChainMain:
		return c.transmitFullChainMain(stream, gson)
	case TopicFullTxPoolAdd:
		return c.transmitFullTxPoolAdd(stream, gson)
	case TopicMinimalChainMain:
		return c.transmitMinimalChainMain(stream, gson)
	case TopicMinimalTxPoolAdd:
		return c.transmitMinimalTxPoolAdd(stream, gson)
	default:
		return fmt.Errorf("unhandled topic '%s'", topic)
	}
}

func (c *Client) transmitFullChainMain(stream *Stream, gson []byte) error {
	arr := []*FullChainMain{}

	if err := json.Unmarshal(gson, &arr); err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}
	for _, element := range arr {
		stream.FullChainMainC <- element
	}

	return nil
}

func (c *Client) transmitFullTxPoolAdd(stream *Stream, gson []byte) error {
	arr := []*FullTxPoolAdd{}

	if err := json.Unmarshal(gson, &arr); err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}
	for _, element := range arr {
		stream.FullTxPoolAddC <- element
	}

	return nil
}

func (c *Client) transmitMinimalChainMain(stream *Stream, gson []byte) error {
	element := &MinimalChainMain{}

	if err := json.Unmarshal(gson, element); err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}
	stream.MinimalChainMainC <- element
	return nil
}

func (c *Client) transmitMinimalTxPoolAdd(stream *Stream, gson []byte) error {
	arr := []*MinimalTxPoolAdd{}

	if err := json.Unmarshal(gson, &arr); err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}
	for _, element := range arr {
		stream.MinimalTxPoolAddC <- element
	}

	return nil
}

func jsonFromFrame(frame []byte) (Topic, []byte, error) {
	unknown := TopicUnknown

	parts := bytes.SplitN(frame, []byte(":"), 2)
	if len(parts) != 2 {
		return unknown, nil, fmt.Errorf(
			"malformed: expected 2 parts, got %d", len(parts))
	}

	topic, gson := string(parts[0]), parts[1]

	switch topic {
	case string(TopicMinimalChainMain):
		return TopicMinimalChainMain, gson, nil
	case string(TopicFullChainMain):
		return TopicFullChainMain, gson, nil
	case string(TopicFullTxPoolAdd):
		return TopicFullTxPoolAdd, gson, nil
	case string(TopicMinimalTxPoolAdd):
		return TopicMinimalTxPoolAdd, gson, nil
	}

	return unknown, nil, fmt.Errorf("unknown topic '%s'", topic)
}
