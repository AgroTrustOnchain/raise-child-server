package entities

import "raise-child/model/dtos/response"

type StaffNft struct {
	ID           ID     `json:"id"`
	IdentityCode string `json:"identity_code"`
	Role         string `json:"role"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Name         string `json:"name"`
	Url          string `json:"url"`
}

func (s StaffNft) ToNftResponse() response.StaffNftResponse {
	if s.ID.ID == "" {
		return response.StaffNftResponse{}
	}

	return response.StaffNftResponse{
		ID:           s.ID.ID,
		IdentityCode: s.IdentityCode,
		Role:         s.Role,
		FirstName:    s.FirstName,
		LastName:     s.LastName,
		Name:         s.Name,
		Url:          s.Url,
	}
}
