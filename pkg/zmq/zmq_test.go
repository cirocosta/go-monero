package zmq_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cirocosta/go-monero/pkg/zmq"
)

func TestJSONFromFrame(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name          string
		input         []byte
		expectedJSON  []byte
		expectedTopic zmq.Topic
		err           string
	}{
		{
			name:  "nil",
			input: nil,
			err:   "malformed",
		},

		{
			name:  "empty",
			input: []byte{},
			err:   "malformed",
		},

		{
			name:  "unknown-topic",
			input: []byte(`foobar:[{"foo":"bar"}]`),
			err:   "unknown topic",
		},

		{
			name:          "proper w/ known-topic",
			input:         []byte(`json-minimal-txpool_add:[{"foo":"bar"}]`),
			expectedTopic: zmq.TopicMinimalTxPoolAdd,
			expectedJSON:  []byte(`[{"foo":"bar"}]`),
		},
	} {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			aTopic, aJSON, err := zmq.JSONFromFrame(tc.input)
			if tc.err != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedTopic, aTopic)
			assert.Equal(t, tc.expectedJSON, aJSON)
		})
	}
}
