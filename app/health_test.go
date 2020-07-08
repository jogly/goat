package app

import (
	"encoding/json"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestHealth(t *testing.T) {
	router := new(gin.Engine)
	fxtest.New(t, Module, fx.Invoke(func(g *gin.Engine) {
		router = g
	}))
	ts := httptest.NewServer(router)
	defer ts.Close()
	client := ts.Client()
	res, err := client.Get(ts.URL + "/health")
	require.NoError(t, err)
	r, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)
	response := gin.H{}
	err = json.Unmarshal(r, &response)
	require.NoError(t, err)

	err = res.Body.Close()
	require.NoError(t, err)
}
