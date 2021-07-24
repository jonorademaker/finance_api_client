package api

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/jonorademaker/finance_api_client/pkg/models"
)

type Api struct {
	OrganisationalAccounts organisationalAccounts
}

type baseApi struct {
	baseUrl       string
	client        *http.Client
	logger        *log.Logger
	retryStrategy retryStrategy
}

type Options struct {
	baseUrl               string
	httpClient            *http.Client
	timeoutInMilliseconds int
	logger                *log.Logger
	retryStrategy         retryStrategy
}

func createBaseApi(options Options) baseApi {
	var logger *log.Logger
	if options.logger != nil {
		logger = options.logger
		logger.SetPrefix("financeapi - ")
	}

	var client *http.Client
	if options.httpClient == nil {
		client = &http.Client{}
	} else {
		client = options.httpClient
	}

	client.Timeout = time.Duration(options.timeoutInMilliseconds) * time.Millisecond

	var baseUrl string
	if options.baseUrl == "" {
		baseUrl = "http://localhost:8080"
	} else {
		baseUrl = options.baseUrl
	}

	var retryStrategy retryStrategy
	if options.retryStrategy != nil {
		retryStrategy = options.retryStrategy
	} else {
		retryStrategy = func(_ RetryOptions) bool { return false }
	}

	return baseApi{baseUrl: baseUrl, client: client, logger: logger, retryStrategy: retryStrategy}
}

func NewApi(options Options) Api {
	api := createBaseApi(options)

	return Api{
		OrganisationalAccounts: organisationalAccounts{baseApi: api},
	}
}

func (api *baseApi) post(resourceUrl string, data []byte) ([]byte, error) {
	fullUrl, err := api.buildUrl(resourceUrl, map[string][]string{})
	if err != nil {
		return nil, err
	}

	resp, err := api.perform("POST", fullUrl, data)
	if err != nil {
		return nil, models.FinanceApiError{Err: err}
	}
	return api.processResponse(resp)
}

func (api *baseApi) get(resourceUrl string) ([]byte, error) {
	fullUrl, err := api.buildUrl(resourceUrl, map[string][]string{})
	if err != nil {
		return nil, err
	}

	resp, err := api.perform("GET", fullUrl, nil)
	if err != nil {
		return nil, models.FinanceApiError{Err: err}
	}
	return api.processResponse(resp)
}

func (api *baseApi) delete(resourceUrl string, queryString map[string][]string) error {
	fullUrl, err := api.buildUrl(resourceUrl, queryString)
	if err != nil {
		return err
	}

	resp, err := api.perform("DELETE", fullUrl, nil)
	if err != nil {
		return models.FinanceApiError{Err: err}
	}

	_, err = api.processResponse(resp)
	return err
}

func (api *baseApi) perform(method string, url *url.URL, data []byte) (*http.Response, error) {
	api.ifLog(func(log *log.Logger) {
		if data != nil {
			log.Printf("%s %s: %s", method, url, string(data))
		} else {
			log.Printf("%s %s", method, url)
		}

	})

	req, err := http.NewRequest(method, url.String(), bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	// required by api docs
	// https://api-docs.form3.tech/api.html#introduction-and-api-conventions-headers
	req.Header.Add("Date", time.Now().Format(time.RFC3339))
	req.Header.Add("User-Agent", "FinanceApi Coding Test Lib v0.0.3")
	req.Header.Add("Accept", "application/vnd.api+json")

	retries := RetryOptions{}
	for {
		retries.response, retries.err = api.client.Do(req)
		retries.count += 1

		if retries.err != nil || !isSuccessResponse(retries.response) {
			api.ifLog(func(log *log.Logger) { log.Printf("response error: %s", retries.err) })
			if api.retryStrategy(retries) {
				continue
			}
		}
		return retries.response, retries.err
	}
}

func (api *baseApi) processResponse(resp *http.Response) ([]byte, error) {
	defer func() {
		if resp != nil && resp.Body != nil {
			err := resp.Body.Close()
			if err != nil {
				api.ifLog(func(log *log.Logger) { log.Printf("failed to close request body: %s", err) })
			}
		}
	}()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, models.FinanceApiError{
			StatusCode: resp.StatusCode,
			Headers:    resp.Header,
			Err:        err,
		}
	}

	if isSuccessResponse(resp) {
		return responseBody, nil
	}

	fullUrl := resp.Request.URL
	err = models.ParseErrorReponse(fullUrl.String(), resp.StatusCode, resp.Header, responseBody)
	return nil, err
}

func (api *baseApi) buildUrl(resourceUrl string, queryString map[string][]string) (*url.URL, error) {
	fullUrl, err := url.Parse(api.baseUrl + resourceUrl)
	if err != nil {
		return nil, models.FinanceApiError{Err: err}
	}

	q := fullUrl.Query()
	for key, value := range queryString {
		for _, v := range value {
			q.Set(key, v)
		}
	}
	fullUrl.RawQuery = q.Encode()
	return fullUrl, nil
}

// ifLog takes a function so we don't do unnecessary work in order to do
// some logging, e.g. decoding bytes to string
func (api *baseApi) ifLog(logFunc func(logger *log.Logger)) {
	if api.logger != nil {
		logFunc(api.logger)
	}
}

func isSuccessResponse(resp *http.Response) bool {
	return resp.StatusCode >= 200 && resp.StatusCode < 300
}
