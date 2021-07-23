#!/usr/bin/env bash

set -ex

ID=`uuidgen`
ORGANISATION_ID=`uuidgen`

curl -v --location --request POST 'http://localhost:8080/v1/organisation/accounts' \
--header 'Date: Fri 11 Jun 06:41:47 BST 2021' \
--data-raw "{
  \"data\": {
    \"id\": \"$ID\",
    \"organisation_id\": \"$ORGANISATION_ID\",
    \"type\": \"accounts\",
    \"attributes\": {
       \"country\": \"GB\",
        \"base_currency\": \"GBP\",
        \"bank_id\": \"400302\",
        \"bank_id_code\": \"GBDSC\",
        \"account_number\": \"10000004\",
        \"customer_id\": \"234\",
        \"iban\": \"GB28NWBK40030212764204\",
        \"bic\": \"NWBKGB42\",
        \"account_classification\": \"Personal\",
        \"name\": [\"John\", \"Doe\"],
        \"alternative_bank_account_names\": [\"123\"],
        \"acceptance_qualifier\": \"123\",
        \"bank_account_name\": \"123\",
        \"first_name\": \"123\",
        \"address\": \"123\",
        \"city\": \"123\",
        \"identification_issuer\": \"123\",
        \"identification_scheme\": \"123\",
        \"registration_number\": \"123\",
        \"tax_residency\": \"123\"
    }
  }
}"
