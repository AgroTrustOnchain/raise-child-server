package entities

import (
	"raise-child/model/dtos/response"
	"raise-child/util"
	"strconv"
)

type BooksNeed struct {
	ID                ID       `json:"id"`
	ChildID           string   `json:"child"`
	Year              string   `json:"year"`
	YearChanges       []string `json:"year_changes"`
	Semster           string   `json:"semester"`
	Value             string   `json:"value"`
	SupportedYears    []string `json:"supported_years"`
	Donors            []string `json:"donors"`
	Donations         []string `json:"donations"`
	WithdrawProposals []string `json:"withdraw_proposals"`
	WithdrawsForNeed  []string `json:"withdraws_for_need"`
	IsUpdated         bool     `json:"is_updated"`
}

type HealthInsuranceNeed struct {
	ID                ID       `json:"id"`
	ChildID           string   `json:"child"`
	Year              string   `json:"year"`
	YearChanges       []string `json:"year_changes"`
	Value             string   `json:"value"`
	SupportedYears    []string `json:"supported_years"`
	Donors            []string `json:"donors"`
	Donations         []string `json:"donations"`
	WithdrawProposals []string `json:"withdraw_proposals"`
	WithdrawsForNeed  []string `json:"withdraws_for_need"`
	IsUpdated         bool     `json:"is_updated"`
}

type MealSupportDuration struct {
	Fields WrapDurationFields `json:"fields"`
}

type WrapVecMap struct {
	Fields VecMapContents `json:"fields"`
}

type VecMapContents struct {
	Contents []VecMapFields `json:"contents"`
}

type VecMapFields struct {
	Fields VecMapWrapFields `json:"fields"`
}

type VecMapWrapFields struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type WrapDurationFields struct {
	StartPeriod string `json:"start_period"`
	EndPeriod   string `json:"end_period"`
}

type MealNeed struct {
	ID                   ID                    `json:"id"`
	ChildID              string                `json:"child"`
	Year                 string                `json:"year"`
	Value                string                `json:"value"`
	Donors               []string              `json:"donors"`
	Donations            []string              `json:"donations"`
	Durations            []MealSupportDuration `json:"durations"`
	TotalSupportedMonths string                `json:"total_supported_months"`
	SupportedYears       WrapVecMap            `json:"supported_years"`
	ProvideMealDates     []string              `json:"provide_meal_dates"`
	ProvideMealPeriods   []string              `json:"provide_meal_periods"`
	ProvideMealStaffs    []string              `json:"provide_meal_staffs"`
	WithdrawProposals    []string              `json:"withdraw_proposals"`
	WithdrawsForNeed     []string              `json:"withdraws_for_need"`
	IsUpdated            bool                  `json:"is_updated"`
}

type OffChainMealSupportDuration struct {
	ID          string
	StartPeriod string
	EndPeriod   string
}

func (m MealNeed) ToMealNeedResponse() response.MealNeedResponse {
	if m.ID.ID == "" {
		return response.MealNeedResponse{}
	}

	year, _ := strconv.Atoi(m.Year)
	value, _ := strconv.ParseInt(m.Value, 10, 64)
	totalMonths, _ := strconv.Atoi(m.TotalSupportedMonths)

	var durations []response.WrapDuration
	var supportedYears []response.WrapSupportedYear
	for i, duration := range m.Durations { // Duration length always longer
		durations = append(durations, response.WrapDuration{
			StartPeriod: util.RawDateToTime(duration.Fields.StartPeriod),
			EndPeriod:   util.RawDateToTime(duration.Fields.EndPeriod),
		})

		if i < len(m.SupportedYears.Fields.Contents) {
			var supportedYear = m.SupportedYears.Fields.Contents[0].Fields
			intYear, _ := strconv.Atoi(supportedYear.Key)
			supportedMonths, _ := strconv.Atoi(supportedYear.Value)

			supportedYears = append(supportedYears, response.WrapSupportedYear{
				Year:            intYear,
				SupportedMonths: supportedMonths,
			})
		}
	}

	return response.MealNeedResponse{
		ID:                   m.ID.ID,
		Year:                 year,
		Value:                value,
		Donors:               m.Donors,
		Donations:            m.Donations,
		Durations:            durations,
		TotalSupportedMonths: totalMonths,
		SupportedYears:       supportedYears,
		WithdrawProposals:    m.WithdrawProposals,
		WithdrawsForNeed:     m.WithdrawsForNeed,
	}
}
