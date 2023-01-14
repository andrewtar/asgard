package bored

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetActivity(t *testing.T) {
	server := setupBoredApi(t, "{\"activity\": \"test_activity\"}", http.StatusOK)
	defer server.Close()

	*boredApiFlag = server.URL + "/api/activity"
	client := BoredClient{}
	response, err := client.GetActivity()

	assert.Nil(t, err)
	assert.Equal(t, response, BoredApiResponse{
		Activity: "test_activity",
	})
}

func TestGetActivityReturnError(t *testing.T) {
	server := setupBoredApi(t, "test_error", http.StatusBadRequest)
	defer server.Close()

	*boredApiFlag = server.URL + "/api/activity"
	client := BoredClient{}
	response, err := client.GetActivity()

	assert.Contains(t, err.Error(), fmt.Sprint(http.StatusBadRequest))
	assert.Equal(t, response, BoredApiResponse{})
}

func setupBoredApi(t *testing.T, response string, status int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/activity" {
			t.Fatalf("Unexpected url, got: %s", r.URL.Path)
			return
		}
		w.WriteHeader(status)
		w.Write([]byte(response))
	}))
}
