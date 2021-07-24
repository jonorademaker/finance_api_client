package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApiRetrySuccess(t *testing.T) {
	count := 1
	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if count == 1 {
			count += 1
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, err := w.Write([]byte("{\"name\": \"John Doe\"}"))
		assert.NoError(t, err)
	}))
	server.Start()
	defer server.Close()

	retry := func(options RetryOptions) bool {
		return (options.response.StatusCode < 200 ||
			options.response.StatusCode >= 300) &&
			options.count < 2
	}
	options := Options{baseUrl: server.URL, httpClient: server.Client(), retryStrategy: retry}
	api := createBaseApi(options)

	resp, err := api.get("/user")
	assert.NoError(t, err)

	body := string(resp)
	assert.Equal(t, "{\"name\": \"John Doe\"}", body)
	assert.Equal(t, 2, count)
}
