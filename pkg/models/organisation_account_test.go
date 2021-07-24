// +build unit

package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOrganisationAccountBasicSerialization(t *testing.T) {
	t.Parallel()
	newAccount := OrganisationAccount{
		ID:             "48e51a61-29e2-44e6-a97d-4bcf3bda92fc",
		OrganisationID: "4f8deb65-3755-4252-a495-9660d00c26a5",
		Attributes: OrganisationAccountAttributes{
			AccountClassification: "Personal",
			Country:               "GB",
			Name:                  []string{"John", "Doe"},
			Status:                "pending",
		},
	}

	payload, err := newAccount.Serialize()
	assert.NoError(t, err)

	expectedPayload := `{
  "data": {
    "type": "accounts",
    "id": "48e51a61-29e2-44e6-a97d-4bcf3bda92fc",
    "organisation_id": "4f8deb65-3755-4252-a495-9660d00c26a5",
    "version": 0,
    "attributes": {
      "country": "GB",
      "customer_id": "",
      "name": [
        "John",
        "Doe"
      ],
      "account_classification": "Personal",
      "joint_account": false,
      "account_matching_opt_out": false,
      "secondary_identification": "",
      "switched": false,
      "status": "pending"
    }
  }
}`

	assert.Equal(t, expectedPayload, string(payload))
}

func TestOrganisationAccountFullSerialization(t *testing.T) {
	t.Parallel()

	location, err := time.LoadLocation("Europe/London")
	assert.NoError(t, err)
	now := time.Date(2021, 6, 13, 19, 51, 27, 813975019, location)

	newAccount := OrganisationAccount{
		ID:             "26628e05-0bbd-4de2-8da4-7d95bcd15ae0",
		OrganisationID: "e13d2e6c-874a-4356-b35a-3e32dab2c34e",
		CreatedOn:      &now,
		ModifiedOn:     &now,
		Attributes: OrganisationAccountAttributes{
			AccountClassification: "Personal",
			AccountNumber:         "10000004",
			AlternativeNames:      []string{"J Dog"},
			BankID:                "400302",
			BankIDCode:            "GBDSC",
			BaseCurrency:          "GBP",
			Bic:                   "NWBKGB42",
			Country:               "GB",
			Iban:                  "GB28NWBK40030212764204",
			Name:                  []string{"John", "Doe"},
			Status:                "pending",
		},
	}

	payload, err := newAccount.Serialize()
	assert.NoError(t, err)

	expectedPayload := `{
  "data": {
    "type": "accounts",
    "id": "26628e05-0bbd-4de2-8da4-7d95bcd15ae0",
    "organisation_id": "e13d2e6c-874a-4356-b35a-3e32dab2c34e",
    "version": 0,
    "created_on": "2021-06-13T19:51:27.813975019+01:00",
    "modified_on": "2021-06-13T19:51:27.813975019+01:00",
    "attributes": {
      "country": "GB",
      "base_currency": "GBP",
      "bank_id": "400302",
      "bank_id_code": "GBDSC",
      "account_number": "10000004",
      "bic": "NWBKGB42",
      "iban": "GB28NWBK40030212764204",
      "customer_id": "",
      "name": [
        "John",
        "Doe"
      ],
      "alternative_names": [
        "J Dog"
      ],
      "account_classification": "Personal",
      "joint_account": false,
      "account_matching_opt_out": false,
      "secondary_identification": "",
      "switched": false,
      "status": "pending"
    }
  }
}`

	assert.Equal(t, expectedPayload, string(payload))
}

func TestOrganisationAccountBasicDeserialization(t *testing.T) {
	payload := `{
  "data": {
    "type": "accounts",
    "id": "48e51a61-29e2-44e6-a97d-4bcf3bda92fc",
    "organisation_id": "4f8deb65-3755-4252-a495-9660d00c26a5",
    "version": 0,
    "attributes": {
      "country": "GB",
      "customer_id": "",
      "name": [
        "John",
        "Doe"
      ],
      "account_classification": "Personal",
      "joint_account": false,
      "account_matching_opt_out": false,
      "secondary_identification": "",
      "switched": false,
      "status": "pending"
    }
  }
}`
	account, err := DeserializeAccountJson([]byte(payload))
	assert.NoError(t, err)

	assert.Equal(t, "48e51a61-29e2-44e6-a97d-4bcf3bda92fc", account.ID)
	assert.Equal(t, "4f8deb65-3755-4252-a495-9660d00c26a5", account.OrganisationID)
	assert.Equal(t, 0, account.Version)
	assert.Nil(t, account.CreatedOn)
	assert.Nil(t, account.ModifiedOn)

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

func TestOrganisationAccountFullDeserialization(t *testing.T) {
	payload := `{
  "data": {
    "type": "accounts",
    "id": "26628e05-0bbd-4de2-8da4-7d95bcd15ae0",
    "organisation_id": "e13d2e6c-874a-4356-b35a-3e32dab2c34e",
    "version": 12,
    "created_on": "2021-06-14T19:51:27.813975019+01:00",
    "modified_on": "2021-06-14T19:51:27.813975019+01:00",
    "attributes": {
      "country": "GB",
      "base_currency": "GBP",
      "bank_id": "400302",
      "bank_id_code": "GBDSC",
      "account_number": "10000004",
      "bic": "NWBKGB42",
      "iban": "GB28NWBK40030212764204",
      "customer_id": "",
      "name": [
        "John",
        "Doe"
      ],
      "alternative_names": [
        "J Dog"
      ],
      "account_classification": "Personal",
      "joint_account": false,
      "account_matching_opt_out": false,
      "secondary_identification": "",
      "switched": false,
      "status": "pending"
    }
  }
}`

	account, err := DeserializeAccountJson([]byte(payload))
	assert.NoError(t, err)

	assert.Equal(t, "26628e05-0bbd-4de2-8da4-7d95bcd15ae0", account.ID)
	assert.Equal(t, "e13d2e6c-874a-4356-b35a-3e32dab2c34e", account.OrganisationID)
	assert.Equal(t, 12, account.Version)

	assert.Equal(t, "2021-06-14T19:51:27.813975019+01:00", account.CreatedOn.Format(time.RFC3339Nano))
	assert.Equal(t, "2021-06-14T19:51:27.813975019+01:00", account.ModifiedOn.Format(time.RFC3339Nano))

	attributes := account.Attributes
	assert.Equal(t, "GB", attributes.Country)
	assert.Equal(t, "GBP", attributes.BaseCurrency)
	assert.Equal(t, "400302", attributes.BankID)
	assert.Equal(t, "GBDSC", attributes.BankIDCode)
	assert.Equal(t, "10000004", attributes.AccountNumber)
	assert.Equal(t, "NWBKGB42", attributes.Bic)
	assert.Equal(t, "GB28NWBK40030212764204", attributes.Iban)
	assert.Equal(t, "", attributes.CustomerId)
	assert.Equal(t, []string{"John", "Doe"}, attributes.Name)
	assert.Equal(t, []string{"J Dog"}, attributes.AlternativeNames)
	assert.Equal(t, "Personal", attributes.AccountClassification)
	assert.Equal(t, false, attributes.JointAccount)
	assert.Equal(t, false, attributes.AccountMatchingOptOut)
	assert.Equal(t, "", attributes.SecondaryIdentification)
	assert.Equal(t, false, attributes.Switched)
	assert.Equal(t, "pending", attributes.Status)
}

func TestOrganisationAccountFailOnIncompleteJson(t *testing.T) {
	payload := `{
  "data": {
    "type": "accounts",
    "id": "26628e05-0bbd-4de2-8da4-7d95bcd15ae0",
    "organisation_id": "e13d2e6c-874a-4356-b35a-3e32dab2c34e",
    "version": 12,
    "created_on": "2021-06-14T19:51:27.813975019+01:00",
    "modified_on": "2021-06-14T19:51:27.813975019+01:00",
    "
}`

	_, err := DeserializeAccountJson([]byte(payload))
	assert.EqualError(t, err, "Error - invalid character '\\n' in string literal")

}
