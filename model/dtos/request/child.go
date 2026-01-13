package request

type GetChildrenRequest struct {
	Keyword     string `json:"keyword"`
	Region      string `json:"region"`
	YearOfBirth *int   `json:"year_of_birth"`
	SortOrder   string `json:"sort_order"`
	Gender      string `json:"gender"`
	Page        int    `json:"page"`
}

type UploadChildRequest struct {
	IdentityCode string `json:"identity_code"`
	FirstName    string `json:"first_name" validate:"required"`
	LastName     string `json:"last_name" validate:"required"`
	Gender       string `json:"gender" validate:"required"`
	DateOfBirth  string `json:"date_of_birth" validate:"required"`
	AvatarBlobId string `json:"avatar_blob_id" validate:"required"`
}

type AddChildStringMetadaRequest struct {
	Key   string `json:"key" validate:"required"`
	Value string `json:"value" validate:"required"`
}

type AddChildNumberMetadaRequest struct {
	Key   string `json:"key" validate:"required"`
	Value int    `json:"value" validate:"required"`
}
