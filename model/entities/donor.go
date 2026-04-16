package entities

import (
	"raise-child/model/dtos/response"
	"strconv"
)

type Donor struct {
	ID              ID       `json:"id"`
	Owner           string   `json:"owner"`
	FirstName       string   `json:"first_name"`
	LastName        string   `json:"last_name"`
	Gender          string   `json:"gender"`
	PhoneNumber     string   `json:"phone_number"`
	Email           string   `json:"email"`
	TotalDonation   string   `json:"total_donation"`
	SupportedChilds []string `json:"supported_childs"`
	Name            string   `json:"name"`
	Url             string   `json:"url"`
}

func (s Donor) ToDonorResponse() response.DonorResponse {
	if s.ID.ID == "" {
		return response.DonorResponse{}
	}

	totalDonation, _ := strconv.ParseInt(s.TotalDonation, 10, 64)

	return response.DonorResponse{
		ID:              s.Owner,
		FirstName:       s.FirstName,
		LastName:        s.LastName,
		Gender:          s.Gender,
		PhoneNumber:     s.PhoneNumber,
		Email:           s.Email,
		TotalDonation:   totalDonation,
		SupportedChilds: s.SupportedChilds,
		Name:            s.Name,
		Url:             s.Url,
	}
}
