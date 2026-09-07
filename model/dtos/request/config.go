package request

type GetPlatformConfigsRequest struct {
	Keyword      string `form:"keyword"`
	ActorAddress string `form:"actor_address"`
	SortOrder    string `form:"sort_order"`
	PageSize     int    `form:"page_size"`
	Page         int    `form:"page"`
}

type UpdatePlatformConfigRequest struct {
	Value       interface{} `json:"value"`
	Description string      `json:"description"`
}
