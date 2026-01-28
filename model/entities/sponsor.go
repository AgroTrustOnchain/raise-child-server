package entities

import (
	"raise-child/model/dtos/response"
	"strconv"
)

type Sponsor struct {
	ID            ID     `json:"id"`
	Owner         string `json:"owner"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Gender        string `json:"gender"`
	PhoneNumber   string `json:"phone_number"`
	Email         string `json:"email"`
	TotalDonation string `json:"total_donation"`
	Name          string `json:"name"`
	Url           string `json:"url"`
}

func (s Sponsor) ToSponsorResponse() response.SponsorResponse {
	if s.ID.ID == "" {
		return response.SponsorResponse{}
	}

	totalDonation, _ := strconv.ParseInt(s.TotalDonation, 10, 64)

	return response.SponsorResponse{
		ID:            s.Owner,
		FirstName:     s.FirstName,
		LastName:      s.LastName,
		Gender:        s.Gender,
		PhoneNumber:   s.PhoneNumber,
		Email:         s.Email,
		TotalDonation: totalDonation,
		Name:          s.Name,
		Url:           s.Url,
	}
}
