package business

import (
	"raise-child/constants/shared"
	"raise-child/interfaces/business"
	"raise-child/model/dtos/response"
	"slices"
)

type regionService struct {
	regions []string
}

var _regions []string

func GenerateRegionService() business.IRegionService {
	if _regions == nil || len(_regions) == 0 {
		_regions = []string{
			shared.TUYEN_QUANG_REGION,
			shared.LAO_CAI_REGION,
			shared.THAI_NGUYEN_REGION,
			shared.PHU_THO_REGION,
			shared.BAC_BINH_REGION,
			shared.HUNG_YEN_REGION,
			shared.HAI_PHONG_REGION,
			shared.NINH_BINH_REGION,
			shared.QUANG_TRI_REGION,
			shared.DA_NANG_REGION,
			shared.QUANG_NGAI_REGION,
			shared.GIA_LAI_REGION,
			shared.KHANH_HOA_REGION,
			shared.LAM_DONG_REGION,
			shared.DAK_LAK_REGION,
			shared.HO_CHI_MINH_REGION,
			shared.DONG_NAI_REGION,
			shared.TAY_NINH_REGION,
			shared.CAN_THO_REGION,
			shared.VINH_LONG_REGION,
			shared.DONG_THAP_REGION,
			shared.CA_MAU_REGION,
			shared.AN_GIANG_REGION,
			shared.HA_NOI_REGION,
			shared.HUE_REGION,
			shared.LAI_CHAU_REGION,
			shared.DIEN_BIEN_REGION,
			shared.SON_LA_REGION,
			shared.QUANG_NINH_REGION,
			shared.THANH_HOA_REGION,
			shared.NGHE_AN_REGION,
			shared.HA_TINH_REGION,
			shared.CAO_BANG_REGION,
		}
	}

	return &regionService{
		regions: _regions,
	}
}

// GetRegions implements business.IRegionService.
func (r *regionService) GetRegions() response.RegionsResponse {
	return response.RegionsResponse{
		Regions: r.regions,
	}
}

func isRegionExist(region string) bool {
	return slices.Contains(_regions, region)
}
