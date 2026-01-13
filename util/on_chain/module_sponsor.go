package onchain

import "raise-child/constants/on-chain/sui"

type IModuleSponsor interface {
	GetModule() string
	GetSponsorNftStruct() string
}

type moduleSponsor struct{}

func InitializeModuleSponsor() IModuleSponsor {
	return &moduleSponsor{}
}

// GetModule implements IModuleSponsor.
func (m *moduleSponsor) GetModule() string {
	return sui.MODULE_SPONSOR
}

// GetSponsorNftStruct implements IModuleSponsor.
func (m *moduleSponsor) GetSponsorNftStruct() string {
	return sui.SPONSOR_NFT_STRUCT
}
