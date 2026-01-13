package entities

import "raise-child/model/dtos/response"

type Sponsor struct {
	ID          ID     `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Gender      string `json:"gender"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	Name        string `json:"name"`
	Url         string `json:"url"`
}

func (s Sponsor) ToSponsorResponse() response.SponsorResponse {
	if s.ID.ID == "" {
		return response.SponsorResponse{}
	}

	return response.SponsorResponse{
		ID:          s.ID.ID,
		FirstName:   s.FirstName,
		LastName:    s.LastName,
		Gender:      s.Gender,
		PhoneNumber: s.PhoneNumber,
		Email:       s.Email,
		Name:        s.Name,
		Url:         s.Url,
	}
}
