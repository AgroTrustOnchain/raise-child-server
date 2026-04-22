package response

type RegionsResponse struct {
	Regions []string `json:"regions"`
}

type RegionDetailResponse struct {
	Region            string                 `json:"region"`
	PoolID            string                 `json:"pool_id"`
	CenterPhoneNumber string                 `json:"center_phone_number"`
	CenterAddress     string                 `json:"center_address"`
	CenterImageBlobID string                 `json:"center_image_blob_id"`
	TotalDonated      int64                  `json:"total_donated"`
	Children          PaginationDataResponse `json:"children"`
}
