package response

import (
	"raise-child/model/dtos/request"
	"time"
)

type ChildResponse struct {
	ID                     string                        `json:"id"`
	IdentityCode           string                        `json:"identity_code"`
	BirthCertificateBlobID string                        `json:"birth_certificate_blob_id"`
	FirstName              string                        `json:"first_name"`
	LastName               string                        `json:"last_name"`
	Gender                 string                        `json:"gender"`
	DateOfBirth            time.Time                     `json:"date_of_birth"`
	HomeAddress            string                        `json:"home_address"`
	Region                 string                        `json:"region"`
	AvatarBlobId           string                        `json:"avatar_blob_id"`
	HomeBlobID             string                        `json:"home_blob_id"`
	FirstGuardian          request.ChildGuardianProfile  `json:"first_guardian"`
	SecondGuardian         *request.ChildGuardianProfile `json:"second_guardian"`
	ImageBlobIds           []string                      `json:"image_blob_ids"`
	UploadImagePeriods     []time.Time                   `json:"upload_image_periods"`
	DynamicFields          []string                      `json:"dynamic_fields"`
	BooksNeeds             []string                      `json:"books_needs"`
	HealthInsuranceNeed    string                        `json:"health_insurance_need"`
	MealNeed               string                        `json:"meal_need"`
	SpecialNeedProposals   []string                      `json:"special_need_proposals"`
	SpecialNeedCampaigns   []string                      `json:"special_need_campaigns"`
	Gifts                  []string                      `json:"gifts"`
	UploadedBy             string                        `json:"uploaded_by"`
	UploadedAt             time.Time                     `json:"uploaded_at"`
	UpdatedAt              time.Time                     `json:"updated_at"`
	DynamicValues          map[string]interface{}        `json:"dynamic_values"`
}

type ChildCardMinimumResponse struct {
	ID           string `json:"id"`
	IdentityCode string `json:"identity_code"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Gender       string `json:"gender"`
}
