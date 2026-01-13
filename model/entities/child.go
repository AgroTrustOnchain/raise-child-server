package entities

import (
	"raise-child/model/dtos/response"
	"raise-child/util"
	"time"
)

type Child struct {
	ID                 ID       `json:"id"`
	IdentityCode       string   `json:"identity_code"`
	FirstName          string   `json:"first_name"`
	LastName           string   `json:"last_name"`
	Gender             string   `json:"gender"`
	DateOfBirth        string   `json:"date_of_birth"`
	Region             string   `json:"region"`
	AvatarBlobId       string   `json:"avatar_blob_id"`
	ImageBlobIds       []string `json:"image_blob_ids"`
	UploadImagePeriods []int64  `json:"upload_image_periods"`
	DynamicFields      []string `json:"dynamic_fields"`
	UploadedAt         int64    `json:"uploaded_at"`
	UpdatedAt          int64    `json:"updated_at"`
}
type ID struct {
	ID string `json:"id"`
}

func (c Child) ToChildResponse() response.ChildResponse {
	if c.ID.ID == "" {
		return response.ChildResponse{}
	}

	var uploadImagePeriods []time.Time
	for _, period := range c.UploadImagePeriods {
		uploadImagePeriods = append(uploadImagePeriods, util.MilliSecToTime(period))
	}

	return response.ChildResponse{
		ID:                 c.ID.ID,
		IdentityCode:       c.IdentityCode,
		FirstName:          c.FirstName,
		LastName:           c.LastName,
		Gender:             c.Gender,
		DateOfBirth:        util.RawDateToTime(c.DateOfBirth),
		Region:             c.Region,
		AvatarBlobId:       c.AvatarBlobId,
		ImageBlobIds:       c.ImageBlobIds,
		UploadImagePeriods: uploadImagePeriods,
		UploadedAt:         util.MilliSecToTime(c.UploadedAt),
		UpdatedAt:          util.MilliSecToTime(c.UpdatedAt),
	}
}
