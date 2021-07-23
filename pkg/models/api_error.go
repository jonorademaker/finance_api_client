package models

import (
	json2 "encoding/json"
	"errors"
	"fmt"
	"regexp"
)

type apiErrorResponse struct {
	ErrorMessage string `json:"error_message"`
}

// called this FinanceApiError because we don't know where
// this error could end up in a user's codebase
type FinanceApiError struct {
	Url        string
	StatusCode int
	Headers    map[string][]string
	RawBody    string
	Err        error
}

func (r FinanceApiError) Error() string {
	if r.StatusCode != 0 {
		return fmt.Sprintf("Error (Status %d) - %v", r.StatusCode, r.Err)
	}

	return fmt.Sprintf("Error - %v", r.Err)
}

func ParseErrorReponse(url string, statusCode int, headers map[string][]string, body []byte) FinanceApiError {
	apiError := FinanceApiError{
		Url:        url,
		StatusCode: statusCode,
		Headers:    headers,
		RawBody:    string(body),
	}

	var apiErrorResponse apiErrorResponse
	err := json2.Unmarshal(body, &apiErrorResponse)
	if err != nil {
		msg := fmt.Sprintf("failed deserialising error response: '%s'", string(body))
		apiError.Err = errors.New(msg)
		return apiError
	}

	sanitizedMessage := removeExtraValidationLines(apiErrorResponse.ErrorMessage)
	apiError.Err = errors.New(sanitizedMessage)
	return apiError
}

func removeExtraValidationLines(errorMessage string) string {
	var re = regexp.MustCompile(`(?m)(validation failure list:\n)+`)
	return re.ReplaceAllString(errorMessage, "validation failure list:\n")
}
