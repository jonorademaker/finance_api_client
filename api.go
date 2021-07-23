package financeapi

import (
	"github.com/jonorademaker/finance_api_client/pkg/api"
)

func NewApi() api.Api {
	// this _should_ paint to the production site
	return api.NewApiWithUrl("http://localhost:8080")
}
