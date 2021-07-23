package api

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/jonorademaker/finance_api_client/pkg/models"
)

type Api struct {
	OrganisationalAccounts organisationalAccounts
}

type baseApi struct {
	baseUrl string
}

func NewApiWithUrl(baseUrl string) Api {
	api := baseApi{baseUrl: baseUrl}

	return Api{
		OrganisationalAccounts: organisationalAccounts{baseApi: api},
	}
}

func (api *baseApi) post(resourceUrl string, data []byte) ([]byte, error) {
	fullUrl, err := api.buildUrl(resourceUrl, map[string][]string{})
	if err != nil {
		return nil, err
	}

	resp, err := perform("POST", fullUrl, data)
	if err != nil {
		return nil, models.FinanceApiError{Err: err}
	}
	return processResponse(resp)
}

func (api *baseApi) get(resourceUrl string) ([]byte, error) {
	fullUrl, err := api.buildUrl(resourceUrl, map[string][]string{})
	if err != nil {
		return nil, err
	}

	resp, err := perform("GET", fullUrl, nil)
	if err != nil {
		return nil, models.FinanceApiError{Err: err}
	}
	return processResponse(resp)
}

func (api *baseApi) delete(resourceUrl string, queryString map[string][]string) error {
	fullUrl, err := api.buildUrl(resourceUrl, queryString)
	if err != nil {
		return err
	}

	resp, err := perform("DELETE", fullUrl, nil)
	if err != nil {
		return models.FinanceApiError{Err: err}
	}

	_, err = processResponse(resp)
	return err
}

func perform(method string, url *url.URL, data []byte) (*http.Response, error) {
	req, err := http.NewRequest(method, url.String(), bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	// required by api docs
	// https://api-docs.form3.tech/api.html#introduction-and-api-conventions-headers
	req.Header.Add("Date", time.Now().Format(time.RFC3339))
	req.Header.Add("User-Agent", "FinanceApi Coding Test Lib v0.0.3")
	req.Header.Add("Accept", "application/vnd.api+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func processResponse(resp *http.Response) ([]byte, error) {
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			fmt.Println("failed to close request body", err)
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

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
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
