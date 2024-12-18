#### Usage

```go
func NewElasticSearchClient(cfg config.Configure, logger logger.Logger) es.ElasticSearchClient {
	var (
		addresses = strings.Split(cfg.GetString("ELASTIC_SEARCH_ADDRESSES"), ",")
		username  = cfg.GetString("ELASTIC_SEARCH_USERNAME")
		password  = cfg.GetString("ELASTIC_SEARCH_PASSWORD")
	)

	client, err := es.NewElasticSearchClient(
		es.WithEsAddresses(addresses),
		es.WithEsUsername(username),
		es.WithEsPassword(password),
	)

	if err != nil {
		panic(err)
	}

	logger.Info(context.TODO(), "Connected to Elasticsearch")
	return client
}

```

```go

package data

import (
	"context"
	"fmt"

	"github.com/lengocson131002/go-clean-core/es"
	"github.com/lengocson131002/go-clean-core/util"
	"github.com/lengocson131002/mcs-account/domain"
	"github.com/lengocson131002/mcs-account/domain/account"
)

const (
	IndexAccountBalance = "t24v2.fbnk_account.transf.1"
)

type accBalanceData struct {
	esClient es.ElasticSearchClient
}

type esAccountBalanceModel struct {
	OcbAccountNumber           string
	OcbBranchCode              string
	OcbCustomerNumber          string
	CustomerNumberJointProfile string
	Currency                   string
	AccountOpeningDate         string
	LastAccountStatusCode      string
	LastAccountStatusDate      string
	Category                   string
	AccountTitle               string
	ShortTitle                 string
	OpenActualBal              int64
	OnlineActualBal            int64
	WorkingBalance             int64
	AccountOfficer             string
	ConditionGroup             string
	CurrNo                     string
	Op_ts                      string
	Current_ts                 string
}

// GetBalance implements data.AccountData.
func (a *accBalanceData) GetBalance(ctx context.Context, accNumber string) (*account.AccountBalanceResponse, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"ocbAccountNumber": accNumber,
			},
		},
	}

	output, err := a.esClient.Search(
		ctx,
		fmt.Sprintf("%s*", IndexAccountBalance),
		es.WithSearchQuery(query),
		es.WithSearchSort([]string{"op_ts:desc"}),
	)

	if err != nil {
		return nil, err
	}

	if len(output) == 0 {
		return nil, domain.ErrorAccountNotFound
	}

	var balRes esAccountBalanceModel
	err = util.MapStruct(output[0], &balRes)
	if err != nil {
		return nil, err
	}

	return &account.AccountBalanceResponse{
		Currency:        balRes.Currency,
		OpenActualBal:   balRes.OpenActualBal,
		WorkingBalance:  balRes.WorkingBalance,
		OnlineActualBal: balRes.OnlineActualBal,
	}, nil
}

func NewAccountBalanceData(esClient es.ElasticSearchClient) account.AccountBalanceData {
	return &accBalanceData{
		esClient: esClient,
	}
}

```