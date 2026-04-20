package entities

import (
	"raise-child/model/dtos/request"
	"raise-child/model/dtos/response"
	"raise-child/util"
	"strconv"
	"time"
)

type Child struct {
	ID                     ID                            `json:"id"`
	IdentityCode           string                        `json:"identity_code"`
	BirthCertificateBlobID string                        `json:"birth_certificate_blob_id"`
	FirstName              string                        `json:"first_name"`
	LastName               string                        `json:"last_name"`
	Gender                 string                        `json:"gender"`
	DateOfBirth            string                        `json:"date_of_birth"`
	HomeAddress            string                        `json:"home_address"`
	Region                 string                        `json:"region"`
	AvatarBlobId           string                        `json:"avatar_blob_id"`
	HomeBlobID             string                        `json:"home_blob_id"`
	GuardianProfiles       []OnChainChildGuardianProfile `json:"guardian_profiles"`
	ImageBlobIds           []string                      `json:"image_blob_ids"`
	UploadImagePeriods     []string                      `json:"upload_image_periods"`
	DynamicFields          []string                      `json:"dynamic_fields"`
	BooksNeeds             []string                      `json:"books_needs"`
	MealNeed               string                        `json:"meal_need"`
	HealthInsuranceNeed    string                        `json:"health_insurance_need"`
	SpecialNeedProposals   []string                      `json:"special_need_proposals"`
	SpecialNeedCampaigns   []string                      `json:"special_need_campaigns"`
	Gifts                  []string                      `json:"gifts"`
	UploadedBy             string                        `json:"uploaded_by"`
	UploadedAt             string                        `json:"uploaded_at"`
	UpdatedAt              string                        `json:"updated_at"`
}

type OnChainChildGuardianProfile struct {
	Fields ChildGuardianProfile `json:"fields"`
}

type ID struct {
	ID string `json:"id"`
}

func (c Child) ToMinimumChildResponse() response.ChildResponse {
	if c.ID.ID == "" {
		return response.ChildResponse{}
	}

	uploadedAt, _ := strconv.ParseInt(c.UploadedAt, 10, 64)
	updatedAt, _ := strconv.ParseInt(c.UpdatedAt, 10, 64)

	var secondGuardianProfile *request.ChildGuardianProfile
	if len(c.GuardianProfiles) > 1 && c.GuardianProfiles[1].Fields.FullName != "" {
		secondGuardianProfile = &request.ChildGuardianProfile{
			FullName:           c.GuardianProfiles[1].Fields.FullName,
			PhoneNumber:        c.GuardianProfiles[1].Fields.PhoneNumber,
			Relation:           c.GuardianProfiles[1].Fields.Relation,
			IdentityCardBlobID: c.GuardianProfiles[1].Fields.IdentityCardBlobID,
		}
	}

	return response.ChildResponse{
		ID:                     c.ID.ID,
		IdentityCode:           c.IdentityCode,
		BirthCertificateBlobID: c.BirthCertificateBlobID,
		FirstName:              c.FirstName,
		LastName:               c.LastName,
		Gender:                 c.Gender,
		DateOfBirth:            util.RawDateToTime(c.DateOfBirth),
		HomeAddress:            c.HomeAddress,
		Region:                 c.Region,
		AvatarBlobId:           c.AvatarBlobId,
		HomeBlobID:             c.HomeBlobID,
		FirstGuardian: request.ChildGuardianProfile{
			FullName:           c.GuardianProfiles[0].Fields.FullName,
			PhoneNumber:        c.GuardianProfiles[0].Fields.PhoneNumber,
			Relation:           c.GuardianProfiles[0].Fields.Relation,
			IdentityCardBlobID: c.GuardianProfiles[0].Fields.IdentityCardBlobID,
		},
		SecondGuardian:       secondGuardianProfile,
		ImageBlobIds:         c.ImageBlobIds,
		BooksNeeds:           c.BooksNeeds,
		HealthInsuranceNeed:  c.HealthInsuranceNeed,
		MealNeed:             c.MealNeed,
		SpecialNeedProposals: c.SpecialNeedProposals,
		SpecialNeedCampaigns: c.SpecialNeedCampaigns,
		Gifts:                c.Gifts,
		UploadedBy:           c.UploadedBy,
		UploadedAt:           util.MilliSecToTime(uploadedAt),
		UpdatedAt:            util.MilliSecToTime(updatedAt),
	}
}

func (c Child) ToChildResponse() response.ChildResponse {
	if c.ID.ID == "" {
		return response.ChildResponse{}
	}

	var uploadImagePeriods []time.Time
	for _, period := range c.UploadImagePeriods {
		uploadPeriod, _ := strconv.ParseInt(period, 10, 64)
		uploadImagePeriods = append(uploadImagePeriods, util.MilliSecToTime(uploadPeriod))
	}

	uploadedAt, _ := strconv.ParseInt(c.UploadedAt, 10, 64)
	updatedAt, _ := strconv.ParseInt(c.UpdatedAt, 10, 64)

	var secondGuardianProfile *request.ChildGuardianProfile
	if len(c.GuardianProfiles) > 1 && c.GuardianProfiles[1].Fields.FullName != "" {
		secondGuardianProfile = &request.ChildGuardianProfile{
			FullName:           c.GuardianProfiles[1].Fields.FullName,
			PhoneNumber:        c.GuardianProfiles[1].Fields.PhoneNumber,
			Relation:           c.GuardianProfiles[1].Fields.Relation,
			IdentityCardBlobID: c.GuardianProfiles[1].Fields.IdentityCardBlobID,
		}
	}

	return response.ChildResponse{
		ID:                     c.ID.ID,
		IdentityCode:           c.IdentityCode,
		BirthCertificateBlobID: c.BirthCertificateBlobID,
		FirstName:              c.FirstName,
		LastName:               c.LastName,
		Gender:                 c.Gender,
		DateOfBirth:            util.RawDateToTime(c.DateOfBirth),
		HomeAddress:            c.HomeAddress,
		Region:                 c.Region,
		AvatarBlobId:           c.AvatarBlobId,
		HomeBlobID:             c.HomeBlobID,
		FirstGuardian: request.ChildGuardianProfile{
			FullName:           c.GuardianProfiles[0].Fields.FullName,
			PhoneNumber:        c.GuardianProfiles[0].Fields.PhoneNumber,
			Relation:           c.GuardianProfiles[0].Fields.Relation,
			IdentityCardBlobID: c.GuardianProfiles[0].Fields.IdentityCardBlobID,
		},
		SecondGuardian:       secondGuardianProfile,
		ImageBlobIds:         c.ImageBlobIds,
		UploadImagePeriods:   uploadImagePeriods,
		BooksNeeds:           c.BooksNeeds,
		HealthInsuranceNeed:  c.HealthInsuranceNeed,
		MealNeed:             c.MealNeed,
		SpecialNeedProposals: c.SpecialNeedProposals,
		SpecialNeedCampaigns: c.SpecialNeedCampaigns,
		Gifts:                c.Gifts,
		UploadedBy:           c.UploadedBy,
		UploadedAt:           util.MilliSecToTime(uploadedAt),
		UpdatedAt:            util.MilliSecToTime(updatedAt),
	}
}

func (c Child) ToChildCardMinimumResponse() response.ChildCardMinimumResponse {
	if c.ID.ID == "" {
		return response.ChildCardMinimumResponse{}
	}

	return response.ChildCardMinimumResponse{
		ID:           c.ID.ID,
		IdentityCode: c.IdentityCode,
		FirstName:    c.FirstName,
		LastName:     c.LastName,
		Gender:       c.Gender,
	}
}
