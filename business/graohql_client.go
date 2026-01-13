package business

import (
	network "raise-child/constants/on-chain/sui"

	"github.com/block-vision/sui-go-sdk/constant"
	"github.com/machinebox/graphql"
)

var _graphQlTestnetClient, _graphQlDevnetClient, _graphQlMainnetClient *graphql.Client
var _graphQlClientAliases map[string]*graphql.Client

func init() {
	if _graphQlTestnetClient == nil {
		_graphQlTestnetClient = graphql.NewClient(network.GRAPHQL_TESTNET_ENDPOINT)
	}

	if _graphQlDevnetClient == nil {
		_graphQlDevnetClient = graphql.NewClient(network.GRAPHQL_DEVNET_ENDPOINT)
	}

	if _graphQlMainnetClient == nil {
		_graphQlMainnetClient = graphql.NewClient(network.GRAPHQL_MAINNET_ENDPOINT)
	}

	_graphQlClientAliases = map[string]*graphql.Client{
		constant.SuiTestnet: _graphQlTestnetClient,
		constant.SuiDevnet:  _graphQlDevnetClient,
		constant.SuiMainnet: _graphQlMainnetClient,
	}
}
