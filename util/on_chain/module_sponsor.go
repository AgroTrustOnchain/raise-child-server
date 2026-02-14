package onchain

import "raise-child/constants/on-chain/sui"

type IModuleDonor interface {
	GetModule() string
	GetDonorNftStruct() string
}

type moduleDonor struct{}

func InitializeModuleDonor() IModuleDonor {
	return &moduleDonor{}
}

// GetModule implements IModuleDonor.
func (m *moduleDonor) GetModule() string {
	return sui.MODULE_DONOR
}

// GetDonorNftStruct implements IModuleDonor.
func (m *moduleDonor) GetDonorNftStruct() string {
	return sui.DONOR_NFT_STRUCT
}
