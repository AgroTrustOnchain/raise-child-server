package onchain

import "raise-child/constants/on-chain/sui"

type IModuleRecord interface {
	GetModule() string
	GetTransactionRecordStruct() string
}

type moduleRecord struct{}

func InitializeModuleRecord() IModuleRecord {
	return &moduleRecord{}
}

// GetModule implements IModuleRecord.
func (m *moduleRecord) GetModule() string {
	return sui.MODULE_RECORD
}

// GetTransactionRecordStruct implements IModuleRecord.
func (m *moduleRecord) GetTransactionRecordStruct() string {
	return sui.TRANSACTION_RECORD_STRUCT
}
