package http_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	mhttp "github.com/cirocosta/go-monero/pkg/http"
)

func TestParseChallenge(t *testing.T) {
	t.Parallel()

	input := `Digest qop="auth",algorithm=MD5,realm="monero-rpc",nonce="IdDHjxbfpLYP/KzjaxaOqA==",stale=false`

	challenge, err := mhttp.ParseChallenge(input)
	assert.NoError(t, err)

	assert.Equal(t, "auth", challenge.Qop)
	assert.Equal(t, "MD5", challenge.Algorithm)
	assert.Equal(t, "monero-rpc", challenge.Realm)
	assert.Equal(t, `IdDHjxbfpLYP/KzjaxaOqA==`, challenge.Nonce)
	assert.Equal(t, "false", challenge.Stale)
}
