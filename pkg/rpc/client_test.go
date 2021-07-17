package rpc_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cirocosta/go-monero/pkg/rpc"
)

func TestClient(t *testing.T) {
	spec.Run(t, "JSONRPC", func(t *testing.T, when spec.G, it spec.S) {
		var (
			ctx    = context.Background()
			client *rpc.Client
			err    error
		)

		it("errors when daemon down", func() {
			daemon := httptest.NewServer(http.HandlerFunc(nil))
			daemon.Close()

			client, err = rpc.NewClient(daemon.URL, rpc.WithHTTPClient(daemon.Client()))
			require.NoError(t, err)

			err = client.JSONRPC(ctx, "method", nil, nil)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "do:")
		})

		it("errors w/ empty response", func() {
			handler := func(w http.ResponseWriter, r *http.Request) {}

			daemon := httptest.NewServer(http.HandlerFunc(handler))
			defer daemon.Close()

			client, err = rpc.NewClient(daemon.URL, rpc.WithHTTPClient(daemon.Client()))
			require.NoError(t, err)

			err = client.JSONRPC(ctx, "method", nil, nil)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "decode")
		})

		it("errors w/ non-200 response", func() {
			handler := func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(500)
			}

			daemon := httptest.NewServer(http.HandlerFunc(handler))
			defer daemon.Close()

			client, err = rpc.NewClient(daemon.URL, rpc.WithHTTPClient(daemon.Client()))
			require.NoError(t, err)

			err = client.JSONRPC(ctx, "method", nil, nil)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "non-2xx status")
		})

		it("makes GET request to the jsonrpc endpoint", func() {
			var (
				endpoint string
				method   string
			)

			handler := func(w http.ResponseWriter, r *http.Request) {
				endpoint = r.URL.Path
				method = r.Method
			}

			daemon := httptest.NewServer(http.HandlerFunc(handler))
			defer daemon.Close()

			client, err = rpc.NewClient(daemon.URL, rpc.WithHTTPClient(daemon.Client()))
			require.NoError(t, err)

			err = client.JSONRPC(ctx, "method", nil, nil)
			assert.Equal(t, rpc.EndpointJSONRPC, endpoint)
			assert.Equal(t, method, "GET")
		})

		it("encodes rpc in request", func() {
			var (
				body = &rpc.RequestEnvelope{}

				params = map[string]interface{}{
					"foo": "bar",
					"caz": 123.123,
				}
			)

			handler := func(w http.ResponseWriter, r *http.Request) {
				err := json.NewDecoder(r.Body).Decode(body)
				assert.NoError(t, err)
			}

			daemon := httptest.NewServer(http.HandlerFunc(handler))
			defer daemon.Close()

			client, err = rpc.NewClient(daemon.URL, rpc.WithHTTPClient(daemon.Client()))
			require.NoError(t, err)

			err = client.JSONRPC(ctx, "rpc-method", params, nil)
			assert.Equal(t, body.ID, "0")
			assert.Equal(t, body.JSONRPC, "2.0")
			assert.Equal(t, body.Method, "rpc-method")
			assert.Equal(t, body.Params, params)
		})

		it("captures result", func() {
			handler := func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, `{"id":"id", "jsonrpc":"jsonrpc", "result": {"foo": "bar"}}`)
			}

			daemon := httptest.NewServer(http.HandlerFunc(handler))
			defer daemon.Close()

			client, err = rpc.NewClient(daemon.URL, rpc.WithHTTPClient(daemon.Client()))
			require.NoError(t, err)

			result := map[string]string{}

			err = client.JSONRPC(ctx, "rpc-method", nil, &result)
			assert.NoError(t, err)

			assert.Equal(t, result, map[string]string{"foo": "bar"})
		})

		it("fails if rpc errored", func() {
			handler := func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, `{"id":"id", "jsonrpc":"jsonrpc", "error": {"code": -1, "message":"foo"}}`)
			}

			daemon := httptest.NewServer(http.HandlerFunc(handler))
			defer daemon.Close()

			client, err = rpc.NewClient(daemon.URL, rpc.WithHTTPClient(daemon.Client()))
			require.NoError(t, err)

			result := map[string]string{}

			err = client.JSONRPC(ctx, "rpc-method", nil, &result)
			assert.Error(t, err)

			assert.Contains(t, err.Error(), "foo")
			assert.Contains(t, err.Error(), "-1")
		})
	}, spec.Report(report.Terminal{}), spec.Parallel(), spec.Random())
}
