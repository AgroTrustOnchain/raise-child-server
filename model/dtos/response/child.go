package response

import "time"

type ChildResponse struct {
	ID                   string                 `json:"id"`
	IdentityCode         string                 `json:"identity_code"`
	FirstName            string                 `json:"first_name"`
	LastName             string                 `json:"last_name"`
	Gender               string                 `json:"gender"`
	DateOfBirth          time.Time              `json:"date_of_birth"`
	Region               string                 `json:"region"`
	AvatarBlobId         string                 `json:"avatar_blob_id"`
	ImageBlobIds         []string               `json:"image_blob_ids"`
	UploadImagePeriods   []time.Time            `json:"upload_image_periods"`
	DynamicFields        []string               `json:"dynamic_fields"`
	BooksNeeds           []string               `json:"books_needs"`
	MealNeed             string                 `json:"meal_need"`
	SpecialNeedProposals []string               `json:"special_need_proposals"`
	SpecialNeedCampaigns []string               `json:"special_need_campaigns"`
	Gifts                []string               `json:"gifts"`
	UploadedBy           string                 `json:"uploaded_by"`
	UploadedAt           time.Time              `json:"uploaded_at"`
	UpdatedAt            time.Time              `json:"updated_at"`
	DynamicValues        map[string]interface{} `json:"dynamic_values"`
}
