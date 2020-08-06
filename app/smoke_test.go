package app

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"

	"github.com/banditml/goat/header"
)

func TestEndpoints(t *testing.T) {
	client, server := makeClientServer(t)
	defer server.Close()

	t.Parallel()
	t.Run("health", func(t *testing.T) {
		res, err := client.Get(server.URL + "/health")
		require.NoError(t, err)
		r, err := ioutil.ReadAll(res.Body)
		require.NoError(t, err)
		response := gin.H{}
		err = json.Unmarshal(r, &response)
		require.NoError(t, err)

		err = res.Body.Close()
		require.NoError(t, err)
	})
}

func makeClientServer(t *testing.T) (*http.Client, *httptest.Server) {
	router := new(gin.Engine)
	app := fxtest.New(t, Module, fx.Populate(&router))
	app.RequireStart()
	server := httptest.NewServer(router)
	client := server.Client()
	withHeaders(client)
	return client, server
}

func withHeaders(c *http.Client) {
	headerTransport := testTransport{
		wrappedTripper: c.Transport,
		headers:        map[string]string{header.BanditID: "test-account"},
	}

	c.Transport = &headerTransport
}

type testTransport struct {
	wrappedTripper http.RoundTripper
	headers        map[string]string
}

// RoundTrip executes a single HTTP transaction, returning
// a Response for the provided Request.
//
// RoundTrip should not attempt to interpret the response. In
// particular, RoundTrip must return err == nil if it obtained
// a response, regardless of the response's HTTP status code.
// A non-nil err should be reserved for failure to obtain a
// response. Similarly, RoundTrip should not attempt to
// handle higher-level protocol details such as redirects,
// authentication, or cookies.
//
// RoundTrip should not modify the request, except for
// consuming and closing the Request's Body. RoundTrip may
// read fields of the request in a separate goroutine. Callers
// should not mutate or reuse the request until the Response's
// Body has been closed.
//
// RoundTrip must always close the body, including on errors,
// but depending on the implementation may do so in a separate
// goroutine even after RoundTrip returns. This means that
// callers wanting to reuse the body for subsequent requests
// must arrange to wait for the Close call before doing so.
//
// The Request's URL and Header fields must be initialized.
func (t *testTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	for key, value := range t.headers {
		r.Header.Add(key, value)
	}
	return t.wrappedTripper.RoundTrip(r)
}
