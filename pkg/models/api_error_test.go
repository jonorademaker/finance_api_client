// +build unit

package models

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupTestVariables(message string) ([]byte, map[string][]string) {
	headers := map[string][]string{
		"Date": {"today"},
	}

	return []byte(message), headers
}

func TestApiErrorDeserializesCorrectly(t *testing.T) {
	t.Parallel()

	errorMessage := "{\"error_message\":\"validation failure list:\\nvalidation failure list:\\nvalidation failure list:\\nstatus in body should be one of [pending confirmed failed]\"}"
	body, headers := setupTestVariables(errorMessage)

	apiError := ParseErrorReponse("/url/test", 400, headers, body)

	expectedError := "Error (Status 400) - validation failure list:\nstatus in body should be one of [pending confirmed failed]"
	assert.EqualError(t, apiError, expectedError)
	assert.Equal(t, "/url/test", apiError.Url)
	assert.Equal(t, 400, apiError.StatusCode)
	assert.Equal(t, errorMessage, apiError.RawBody)
	assert.Equal(t, headers, headers)
}

func TestApiErrorDeserializationHandlesMarhsalError(t *testing.T) {
	t.Parallel()
	errorMessage := "{\"error_message\": \"invalid message}"
	body, headers := setupTestVariables(errorMessage)

	apiError := ParseErrorReponse("/url/test2", 500, headers, body)

	expectedError := "Error (Status 500) - failed deserialising error response: '{\"error_message\": \"invalid message}'"
	assert.EqualError(t, apiError, expectedError)
	assert.Equal(t, "/url/test2", apiError.Url)
	assert.Equal(t, 500, apiError.StatusCode)
	assert.Equal(t, errorMessage, apiError.RawBody)
	assert.Equal(t, headers, headers)
}

func TestApiErrorHandleNoStatusCode(t *testing.T) {
	t.Parallel()

	err := FinanceApiError{Err: errors.New("couldn't access localhost")}
	assert.EqualError(t, err, "Error - couldn't access localhost")

}
