package daemon

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cirocosta/go-monero/cmd/monero/display"
	"github.com/cirocosta/go-monero/pkg/zmq"
)

type zmqCommand struct {
	JSON bool

	topic    string
	endpoint string
}

var zmqTopics = []string{
	string(zmq.TopicMinimalTxPoolAdd),
	string(zmq.TopicFullTxPoolAdd),
	string(zmq.TopicMinimalChainMain),
	string(zmq.TopicFullChainMain),
}

func (c *zmqCommand) Cmd() *cobra.Command {
	var topicChoicesTxt = fmt.Sprintf("(%s)", strings.Join(zmqTopics, ","))

	cmd := &cobra.Command{
		Use:   "zmq",
		Short: "listen for zmq notifications",
		RunE:  c.RunE,
	}

	cmd.Flags().BoolVar(&c.JSON, "json",
		false, "whether or not to output the result as json")

	cmd.Flags().StringVar(&c.topic, "topic",
		"json-minimal-txpool_add", "zmq topic to subscribe to "+
			topicChoicesTxt)
	_ = cmd.MarkFlagRequired("topic")
	_ = cmd.RegisterFlagCompletionFunc("topic", c.topicCompletion)

	cmd.Flags().StringVar(&c.endpoint, "endpoint",
		"", "zero-mq endpoint to listen for publications")
	_ = cmd.MarkFlagRequired("endpoint")

	return cmd
}

func (c *zmqCommand) topicCompletion(
	cmd *cobra.Command, args []string, toComplete string,
) ([]string, cobra.ShellCompDirective) {
	return zmqTopics, cobra.ShellCompDirectiveDefault
}

func (c *zmqCommand) RunE(_ *cobra.Command, _ []string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client := zmq.NewClient(c.endpoint, zmq.Topic(c.topic))
	defer client.Close()

	stream, err := client.Listen(ctx)
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}

	if err := c.consumeStream(ctx, stream); err != nil {
		return fmt.Errorf("consume stream: %w", err)
	}

	return nil
}

func (c *zmqCommand) consumeStream(
	ctx context.Context, stream *zmq.Stream,
) error {
	for {
		var tx interface{}

		select {
		case err := <-stream.ErrC:
			return err
		case <-ctx.Done():
			return ctx.Err()
		case tx = <-stream.FullChainMainC:
		case tx = <-stream.FullTxPoolAddC:
		case tx = <-stream.MinimalChainMainC:
		case tx = <-stream.MinimalTxPoolAddC:
		}

		if tx != nil {
			if err := display.JSON(tx); err != nil {
				return fmt.Errorf("display json: %w", err)
			}
		}
	}
}

func init() {
	RootCommand.AddCommand((&zmqCommand{}).Cmd())
}
