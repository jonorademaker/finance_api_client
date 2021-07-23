// +build integration

package api

import (
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/jonorademaker/finance_api_client/pkg/models"
	"github.com/stretchr/testify/assert"
)

func createTestApi() Api {
	testUrl := os.Getenv("TEST_URL")
	if testUrl == "" {
		return NewApiWithUrl("http://localhost:8080")
	}

	return NewApiWithUrl(testUrl)
}

func createAccount(t *testing.T) models.OrganisationAccount {
	accountId, err := uuid.NewUUID()
	assert.NoError(t, err)

	organisationId, err := uuid.NewUUID()
	assert.NoError(t, err)

	newAccount := models.OrganisationAccount{
		ID:             accountId.String(),
		OrganisationID: organisationId.String(),
		Attributes: models.OrganisationAccountAttributes{
			AccountClassification: "Personal",
			Country:               "GB",
			Name:                  []string{"John", "Doe"},
			Status:                "pending",
		},
	}
	return newAccount
}

func TestApiOrganisationAccountCreate(t *testing.T) {
	newAccount := createAccount(t)
	api := createTestApi()

	account, err := api.OrganisationalAccounts.Create(newAccount)

	assert.NoError(t, err)

	assert.Equal(t, newAccount.ID, account.ID)
	assert.Equal(t, newAccount.OrganisationID, account.OrganisationID)
	assert.Equal(t, 0, account.Version)
	assert.NotNil(t, account.CreatedOn)
	assert.NotNil(t, account.ModifiedOn)

	attributes := account.Attributes
	assert.Equal(t, "GB", attributes.Country)
	assert.Equal(t, "", attributes.CustomerId)
	assert.Equal(t, []string{"John", "Doe"}, attributes.Name)
	assert.Equal(t, "Personal", attributes.AccountClassification)
	assert.Equal(t, false, attributes.JointAccount)
	assert.Equal(t, false, attributes.AccountMatchingOptOut)
	assert.Equal(t, "", attributes.SecondaryIdentification)
	assert.Equal(t, false, attributes.Switched)
	assert.Equal(t, "pending", attributes.Status)
}

func TestApiOrganisationAccountCreateFail(t *testing.T) {
	newAccount := createAccount(t)
	newAccount.ID = ""
	newAccount.OrganisationID = ""
	api := createTestApi()

	_, err := api.OrganisationalAccounts.Create(newAccount)

	assert.EqualError(t, err, "Error (Status 400) - validation failure list:\nid in body is required\norganisation_id in body is required")
}

func TestApiOrganisationAccountFetch(t *testing.T) {
	newAccount := createAccount(t)
	api := createTestApi()
	account, err := api.OrganisationalAccounts.Create(newAccount)
	assert.NoError(t, err)

	fetchedAccount, err := api.OrganisationalAccounts.Fetch(newAccount.ID)
	assert.NoError(t, err)

	assert.Equal(t, newAccount.ID, fetchedAccount.ID)
	assert.Equal(t, newAccount.OrganisationID, fetchedAccount.OrganisationID)
	assert.Equal(t, newAccount.Version, fetchedAccount.Version)
	assert.Equal(t, account.CreatedOn, fetchedAccount.CreatedOn)
	assert.Equal(t, account.ModifiedOn, fetchedAccount.ModifiedOn)

	newAttributes := newAccount.Attributes
	attributes := fetchedAccount.Attributes
	assert.Equal(t, newAttributes.Country, attributes.Country)
	assert.Equal(t, newAttributes.CustomerId, attributes.CustomerId)
	assert.Equal(t, newAttributes.Name, attributes.Name)
	assert.Equal(t, newAttributes.AccountClassification, attributes.AccountClassification)
	assert.Equal(t, newAttributes.JointAccount, attributes.JointAccount)
	assert.Equal(t, newAttributes.AccountMatchingOptOut, attributes.AccountMatchingOptOut)
	assert.Equal(t, newAttributes.SecondaryIdentification, attributes.SecondaryIdentification)
	assert.Equal(t, newAttributes.Switched, attributes.Switched)
	assert.Equal(t, newAttributes.Status, attributes.Status)
}

func TestApiOrganisationAccountFetchFail(t *testing.T) {
	api := createTestApi()

	_, err := api.OrganisationalAccounts.Fetch("fab96ee7-edb3-4319-b486-a1233ad52960")

	assert.EqualError(t, err, "Error (Status 404) - record fab96ee7-edb3-4319-b486-a1233ad52960 does not exist")
}

func TestApiOrganisationAccountDelete(t *testing.T) {
	newAccount := createAccount(t)
	api := createTestApi()
	_, err := api.OrganisationalAccounts.Create(newAccount)
	assert.NoError(t, err)

	err = api.OrganisationalAccounts.Delete(newAccount.ID, 0)
	assert.NoError(t, err)

	_, err = api.OrganisationalAccounts.Fetch(newAccount.ID)
	assert.Errorf(t, err, "Error (Status 404) - record %s does not exist", newAccount.ID)
}

func TestApiOrganisationAccountDeleteFail(t *testing.T) {
	api := createTestApi()

	err := api.OrganisationalAccounts.Delete("adsfas7fasfdsdaf", 234)
	assert.EqualError(t, err, "Error (Status 400) - id is not a valid uuid")
}
