package repository

import (
	"context"
	"raise-child/model/dtos/request"
	"raise-child/model/entities"
)

type IPlatformConfigRepository interface {
	GetConfigs(req request.GetPlatformConfigsRequest, table string, ctx context.Context) ([]entities.PlatformConfig, int, error)
	GetConfig(id, table string, ctx context.Context) (*entities.PlatformConfig, error)
	GetConfigByKey(table, key string, ctx context.Context) (*entities.PlatformConfig, error)
	CreateConfig(config entities.PlatformConfig, table string, ctx context.Context) error
	UpdateConfig(config entities.PlatformConfig, table string, ctx context.Context) error
}
