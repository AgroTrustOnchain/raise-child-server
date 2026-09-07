package business

import (
	network "raise-child/constants/on-chain/sui"

	"github.com/block-vision/sui-go-sdk/constant"
	"github.com/block-vision/sui-go-sdk/sui"
)

var _testnetClient, _devnetClient, _mainnetClient sui.ISuiAPI
var _networkAliases map[string]sui.ISuiAPI

func init() {
	if _testnetClient == nil {
		_testnetClient = sui.NewSuiClient(network.TESTNET_NETWORK) // Optimized SUI testnet gRPC url from Block Vision (Third Party)
	}

	if _devnetClient == nil {
		_devnetClient = sui.NewSuiClient(network.DEVNET_NETWORK) // Official SUI devnet gRPC url (Block Vision does not support this network)
	}

	if _mainnetClient == nil {
		_mainnetClient = sui.NewSuiClient(constant.BvMainnetEndpoint) // Optimized SUI mainnet gRPC url from Block Vision (Third Party)
	}

	_networkAliases = map[string]sui.ISuiAPI{
		constant.SuiTestnet: _testnetClient,
		constant.SuiDevnet:  _devnetClient,
		constant.SuiMainnet: _mainnetClient,
	}
}
