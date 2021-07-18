//
// This is a derived work based on `code.google.com/p/mlab-ns2/gae/ns/digest`
// (original work of Bipasa Chattopadhyay bipasa@cs.unc.edu Eric Gavaletz
// gavaletz@gmail.com Seon-Wook Park seon.wook@swook.net, from the fork
// maintained by Bob Ziuchkovski @bobziuchkovski
// (https://github.com/rkl-/digest).
//
package http

import (
	"bytes"
	"crypto/md5" // nolint:gosec
	"crypto/rand"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

// DigestAuthTransport is an implementation of http.RoundTripper that takes
// care of http digest authentication.
//
type DigestAuthTransport struct {
	Username string
	Password string
	rt       http.RoundTripper
}

// NewDigestAuthTransport creates a new digest transport using the
// http.DefaultTransport.
//
func NewDigestAuthTransport(
	username, password string, rt http.RoundTripper,
) *DigestAuthTransport {
	return &DigestAuthTransport{
		Username: username,
		Password: password,
		rt:       rt,
	}
}

func (t *DigestAuthTransport) newCredentials(
	req *http.Request, c *challenge,
) *credentials {
	return &credentials{
		Algorithm:  c.Algorithm,
		DigestURI:  req.URL.RequestURI(),
		MessageQop: c.Qop, // "auth" must be a single value
		Nonce:      c.Nonce,
		NonceCount: 0,
		Opaque:     c.Opaque,
		Realm:      c.Realm,
		Username:   t.Username,

		method:   req.Method,
		password: t.Password,
	}
}

// RoundTrip makes a request expecting a 401 response that will require digest
// authentication. It creates the credentials it needs and makes a follow-up
// request.
//
func (t *DigestAuthTransport) RoundTrip(
	req *http.Request,
) (*http.Response, error) {
	finalRequest := req.Clone(req.Context())

	if req.Body != nil {
		bodyContents, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return nil, fmt.Errorf("read all body: %w", err)
		}

		reqBody01 := io.NopCloser(bytes.NewBuffer(bodyContents))
		reqBody02 := io.NopCloser(bytes.NewBuffer(bodyContents))

		req.Body = reqBody01
		finalRequest.Body = reqBody02
	}

	// make a request to get the 401 that contains the challenge.
	//
	resp, err := t.rt.RoundTrip(req)
	if err != nil {
		return nil, fmt.Errorf("round trip err: %w", err)
	}
	if resp.StatusCode != 401 { // cool, reached what we needed
		return resp, nil
	}

	// we must ensure that the initial response has been totally drained
	// otherwise the http client won't reuse the connection.
	//
	if _, err := io.Copy(ioutil.Discard, resp.Body); err != nil {
		return nil, fmt.Errorf("copy body to null dev: %w", err)
	}
	resp.Body.Close()

	chal, err := parseChallenge(resp.Header.Get("WWW-Authenticate"))
	if err != nil {
		return nil, fmt.Errorf("parse challenge: %w", err)
	}

	cr := t.newCredentials(finalRequest, chal)
	auth, err := cr.authorize()
	if err != nil {
		return nil, fmt.Errorf("authorize: %w", err)
	}

	finalRequest.Header.Set("Authorization", auth)
	return t.rt.RoundTrip(finalRequest)
}

type challenge struct {
	Realm     string
	Domain    string
	Nonce     string
	Opaque    string
	Stale     string
	Algorithm string
	Qop       string
}

func parseChallenge(input string) (*challenge, error) {
	const quotation = `"`

	fields, err := parseChallengeFields(input)
	if err != nil {
		return nil, fmt.Errorf("parse challenge fields: %w", err)
	}

	c := &challenge{}
	for _, field := range fields {
		kv := strings.SplitN(field, "=", 2)
		if len(kv) != 2 {
			return nil, fmt.Errorf("split: expected to 2 parts "+
				"in entry, got %d. field: '%s'",
				len(kv), field)
		}

		key, value := kv[0], kv[1]
		value = strings.Trim(value, quotation)
		switch key {
		case "qop":
			c.Qop = value
		case "algorithm":
			c.Algorithm = value
		case "realm":
			c.Realm = value
		case "nonce":
			c.Nonce = value
		case "stale":
			c.Stale = value
		default:
			return nil, fmt.Errorf("unknown field '%s'", key)
		}
	}

	return c, nil
}

func parseChallengeFields(str string) ([]string, error) {
	const challengePrefix = "Digest "
	const whitespaceDelimiters = " \n\r\t"

	str = strings.Trim(str, whitespaceDelimiters)
	if !strings.HasPrefix(str, challengePrefix) {
		return nil, fmt.Errorf("bad challenge: "+
			"input doesn't start with '%s'", challengePrefix)
	}

	str = strings.Trim(str[len(challengePrefix):], whitespaceDelimiters)
	fields := strings.Split(str, ",")
	if len(fields) != 5 {
		return nil, fmt.Errorf("split: expected 5 fields, got %d",
			len(fields))
	}

	return fields, nil
}

type credentials struct {
	Algorithm  string
	Cnonce     string
	DigestURI  string
	MessageQop string
	Nonce      string
	NonceCount int
	Opaque     string
	Realm      string
	Username   string

	method   string
	password string
}

func (c *credentials) ha1() string {
	return h(fmt.Sprintf("%s:%s:%s", c.Username, c.Realm, c.password))
}

func (c *credentials) ha2() string {
	return h(fmt.Sprintf("%s:%s", c.method, c.DigestURI))
}

func (c *credentials) resp() (string, error) {
	c.NonceCount++

	if c.MessageQop != "auth" {
		return "", fmt.Errorf("unexpected messageqop '%s'",
			c.MessageQop)
	}

	b := make([]byte, 8)
	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		return "", fmt.Errorf("read full: %w", err)
	}

	c.Cnonce = fmt.Sprintf("%x", b)[:16]

	data := fmt.Sprintf("%s:%08x:%s:%s:%s",
		c.Nonce, c.NonceCount, c.Cnonce, c.MessageQop, c.ha2())
	return kd(c.ha1(), data), nil
}

func (c *credentials) authorize() (string, error) {
	// Note that this is only implemented for MD5 and NOT MD5-sess.
	// MD5-sess is rarely supported and those that do are a big mess.
	if c.Algorithm != "MD5" {
		return "", fmt.Errorf("unsupported algorithm '%s'",
			c.Algorithm)
	}

	resp, err := c.resp()
	if err != nil {
		return "", fmt.Errorf("resp: %w", err)
	}

	sl := []string{fmt.Sprintf(`username="%s"`, c.Username)}
	sl = append(sl, fmt.Sprintf(`realm="%s"`, c.Realm))
	sl = append(sl, fmt.Sprintf(`nonce="%s"`, c.Nonce))
	sl = append(sl, fmt.Sprintf(`uri="%s"`, c.DigestURI))
	sl = append(sl, fmt.Sprintf(`response="%s"`, resp))

	if c.Algorithm != "" {
		sl = append(sl, fmt.Sprintf(`algorithm="%s"`, c.Algorithm))
	}

	if c.Opaque != "" {
		sl = append(sl, fmt.Sprintf(`opaque="%s"`, c.Opaque))
	}

	if c.MessageQop != "" {
		sl = append(sl, fmt.Sprintf("qop=%s", c.MessageQop))
		sl = append(sl, fmt.Sprintf("nc=%08x", c.NonceCount))
		sl = append(sl, fmt.Sprintf(`cnonce="%s"`, c.Cnonce))
	}

	return fmt.Sprintf("Digest %s", strings.Join(sl, ", ")), nil
}

func h(data string) string {
	// `gosec` won't be happy ("weak crypto primitive"), but it's what the
	// server uses.
	//
	// nolint:gosec
	hf := md5.New()
	if _, err := io.WriteString(hf, data); err != nil {
		panic(fmt.Errorf("write string: %w", err))
	}
	return fmt.Sprintf("%x", hf.Sum(nil))
}

func kd(secret, data string) string {
	return h(fmt.Sprintf("%s:%s", secret, data))
}
