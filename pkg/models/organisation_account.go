package models

import (
	json2 "encoding/json"
	"time"
)

type accountCreation struct {
	Data OrganisationAccount `json:"data"`
}

type OrganisationAccountAttributes struct {
	Country                 string   `json:"country"`
	BaseCurrency            string   `json:"base_currency,omitempty"`
	BankID                  string   `json:"bank_id,omitempty"`
	BankIDCode              string   `json:"bank_id_code,omitempty"`
	AccountNumber           string   `json:"account_number,omitempty"`
	Bic                     string   `json:"bic,omitempty"`
	Iban                    string   `json:"iban,omitempty"`
	CustomerId              string   `json:"customer_id"`
	Name                    []string `json:"name"`
	AlternativeNames        []string `json:"alternative_names,omitempty"`
	AccountClassification   string   `json:"account_classification"`
	JointAccount            bool     `json:"joint_account"`
	AccountMatchingOptOut   bool     `json:"account_matching_opt_out"`
	SecondaryIdentification string   `json:"secondary_identification"`
	Switched                bool     `json:"switched"`
	Status                  string   `json:"status"`
}

type OrganisationAccount struct {
	Type           string                        `json:"type"`
	ID             string                        `json:"id"`
	OrganisationID string                        `json:"organisation_id"`
	Version        int                           `json:"version"`
	CreatedOn      *time.Time                    `json:"created_on,omitempty"`
	ModifiedOn     *time.Time                    `json:"modified_on,omitempty"`
	Attributes     OrganisationAccountAttributes `json:"attributes"`
}

func (account *OrganisationAccount) Serialize() ([]byte, error) {
	newAccount := accountCreation{Data: *account}
	newAccount.Data.Type = "accounts"
	data, err := json2.MarshalIndent(newAccount, "", "  ")
	if err != nil {
		return nil, FinanceApiError{Err: err}
	}
	return data, nil
}

func DeserializeAccountJson(body []byte) (OrganisationAccount, error) {
	var savedAccount accountCreation
	err := json2.Unmarshal(body, &savedAccount)
	if err != nil {
		return OrganisationAccount{}, FinanceApiError{Err: err}
	}

	return savedAccount.Data, nil
}
