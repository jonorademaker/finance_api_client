package financeapi

import (
	"github.com/jonorademaker/finance_api_client/pkg/api"
)

func NewApi() api.Api {
	// this _should_ point to the production site
	return api.NewApi(api.Options{})
}
