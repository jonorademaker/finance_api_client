package api

import (
	"fmt"
	"strconv"

	"github.com/jonorademaker/finance_api_client/pkg/models"
)

type organisationalAccounts struct {
	baseApi baseApi
}

func (accountsApi *organisationalAccounts) Create(account models.OrganisationAccount) (models.OrganisationAccount, error) {
	payload, err := account.Serialize()
	if err != nil {
		return models.OrganisationAccount{}, err
	}

	responseBody, err := accountsApi.baseApi.post("/v1/organisation/accounts", payload)
	if err != nil {
		return models.OrganisationAccount{}, err
	}

	return models.DeserializeAccountJson(responseBody)
}

func (accountsApi *organisationalAccounts) Fetch(id string) (models.OrganisationAccount, error) {
	resourceUrl := fmt.Sprintf("/v1/organisation/accounts/%s", id)
	responseBody, err := accountsApi.baseApi.get(resourceUrl)
	if err != nil {
		return models.OrganisationAccount{}, err
	}

	return models.DeserializeAccountJson(responseBody)
}

func (accountsApi *organisationalAccounts) Delete(id string, version int) error {
	resourceUrl := fmt.Sprintf("/v1/organisation/accounts/%s", id)
	queryString := map[string][]string{"version": {strconv.Itoa(version)}}

	return accountsApi.baseApi.delete(resourceUrl, queryString)
}
