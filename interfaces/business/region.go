package business

import (
	"raise-child/model/dtos/response"
)

type IRegionService interface {
	GetRegions() response.RegionsResponse
}
