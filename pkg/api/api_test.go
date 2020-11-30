package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/finitum/aurum/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestGetPublicKey(t *testing.T) {
	pkrsp := models.PublicKeyResponse{PublicKey: "apublickey"}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/pk", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		err := json.NewEncoder(w).Encode(&pkrsp)
		assert.NoError(t, err)
	}))
	defer ts.Close()

	resp, err := GetPublicKey(ts.URL)
	assert.NoError(t, err)

	assert.Equal(t, &pkrsp, resp)
}
