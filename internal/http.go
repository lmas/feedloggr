package internal

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"
)

func init() {
	// It's such a low impact use case, there's no need for crypto/rand.
	// Famous last words...
	rand.Seed(time.Now().Unix())
}

const millis int = 1000
const minRate int = 500
const defaultJitter int = 2

// RateLimit accepts a jitter value (in seconds), uses it to create a new random
// timeout (in milliseconds) and adds a minimum rate (500ms).
// If jitter is less than one second, a default value will be used (2s).
func rateLimit(jitter int) time.Duration {
	if jitter < 1 {
		jitter = defaultJitter
	}
	// gosec warns about bad rand source, but see above note for init()
	i := rand.Intn(jitter*millis) + minRate // #nosec G404
	return time.Duration(i) * time.Millisecond
}

////////////////////////////////////////////////////////////////////////////////

type transport struct {
	http.RoundTripper
}

func newTransport(dir string) *transport {
	d := http.DefaultTransport
	// The file protocol enables easier testing
	d.(*http.Transport).RegisterProtocol("file", http.NewFileTransport(http.Dir(dir)))
	return &transport{d}
}

func (t *transport) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Set("User-Agent", "linux:"+GeneratorName+":"+GeneratorVersion+" ("+GeneratorSource+")")
	return t.RoundTripper.RoundTrip(r)
}

type clientConf struct {
	Timeout int
	Jitter  int
}

type client struct {
	http *http.Client
	conf clientConf
}

func newClient(conf clientConf) *client {
	h := &http.Client{
		Transport: newTransport("."),
		Timeout:   time.Duration(conf.Timeout) * time.Second,
	}
	return &client{
		http: h,
		conf: conf,
	}
}

func (c *client) Get(path string) (io.ReadCloser, error) {
	r, err := c.http.Get(path)
	if err != nil {
		return nil, err
	}
	if r.StatusCode < 200 || r.StatusCode > 299 {
		// NOTE: program run shouldn't be long lived (like a server) so should
		// be safe to ignore any .Close() errors
		// gosec warns about unhandled errors, but we don't really care about that here
		r.Body.Close() //#nosec G104
		return nil, fmt.Errorf("bad response status: %s", r.Status)
	}
	return r.Body, nil
}

func (c *client) RateLimitedGet(path string) (io.ReadCloser, error) {
	d := rateLimit(c.conf.Jitter)
	time.Sleep(d)
	return c.Get(path)
}
