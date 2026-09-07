package entities

import (
	"raise-child/model/dtos/response"
	"raise-child/util"
	"strconv"
)

type BooksNeed struct {
	ID                      ID       `json:"id"`
	ChildID                 string   `json:"child"`
	Year                    string   `json:"year"`
	YearChanges             []string `json:"year_changes"`
	ValueChanges            []string `json:"value_changes"`
	Semster                 string   `json:"semester"`
	Value                   string   `json:"value"`
	SupportedYears          []string `json:"supported_years"`
	Donors                  []string `json:"donors"`
	DonorTotalContributions []string `json:"donor_total_contributions"`
	ProvideDates            []string `json:"provide_dates"`
	ProvidePeriods          []string `json:"provide_periods"`
	ProvideStaffs           []string `json:"provide_staffs"`
	ProvideImageBlobIDs     []string `json:"provide_image_blob_ids"`
	Donations               []string `json:"donations"`
	WithdrawProposals       []string `json:"withdraw_proposals"`
	WithdrawsForNeed        []string `json:"withdraws_for_need"`
	IsUpdated               bool     `json:"is_updated"`
}

type HealthInsuranceNeed struct {
	ID                      ID       `json:"id"`
	ChildID                 string   `json:"child"`
	Year                    string   `json:"year"`
	YearChanges             []string `json:"year_changes"`
	ValueChanges            []string `json:"value_changes"`
	Value                   string   `json:"value"`
	SupportedYears          []string `json:"supported_years"`
	Donors                  []string `json:"donors"`
	DonorTotalContributions []string `json:"donor_total_contributions"`
	ProvideDates            []string `json:"provide_dates"`
	ProvidePeriods          []string `json:"provide_periods"`
	ProvideStaffs           []string `json:"provide_staffs"`
	ProvideImageBlobIDs     []string `json:"provide_image_blob_ids"`
	Donations               []string `json:"donations"`
	WithdrawProposals       []string `json:"withdraw_proposals"`
	WithdrawsForNeed        []string `json:"withdraws_for_need"`
	IsUpdated               bool     `json:"is_updated"`
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
	ID                      ID                    `json:"id"`
	ChildID                 string                `json:"child"`
	Year                    string                `json:"year"`
	Value                   string                `json:"value"`
	YearChanges             []string              `json:"year_changes"`
	ValueChanges            []string              `json:"value_changes"`
	Donors                  []string              `json:"donors"`
	DonorTotalContributions []string              `json:"donor_total_contributions"`
	Donations               []string              `json:"donations"`
	Durations               []MealSupportDuration `json:"durations"`
	TotalSupportedMonths    string                `json:"total_supported_months"`
	SupportedYears          WrapVecMap            `json:"supported_years"`
	ProvideDates            []string              `json:"provide_dates"`
	ProvidePeriods          []string              `json:"provide_periods"`
	ProvideStaffs           []string              `json:"provide_staffs"`
	ProvideImageBlobIDs     []string              `json:"provide_image_blob_ids"`
	WithdrawProposals       []string              `json:"withdraw_proposals"`
	WithdrawsForNeed        []string              `json:"withdraws_for_need"`
	IsUpdated               bool                  `json:"is_updated"`
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

	// var provideMealPeriods []time.Time
	// for _, rawPeriod := range m.ProvidePeriods {
	// 	period, _ := strconv.ParseInt(rawPeriod, 10, 64)
	// 	provideMealPeriods = append(provideMealPeriods, util.MilliSecToTime(period))
	// }

	return response.MealNeedResponse{
		ID:                   m.ID.ID,
		Year:                 year,
		Value:                value,
		Donors:               m.Donors,
		Donations:            m.Donations,
		Durations:            durations,
		TotalSupportedMonths: totalMonths,
		SupportedYears:       supportedYears,
		ProvideDates:         m.ProvideDates,
		ProvidePeriods:       m.ProvidePeriods,
		ProvideStaffs:        m.ProvideStaffs,
		ProvideImageBlobIDs:  m.ProvideImageBlobIDs,
		WithdrawProposals:    m.WithdrawProposals,
		WithdrawsForNeed:     m.WithdrawsForNeed,
	}
}

func (b BooksNeed) ToBooksNeedReponse() response.BooksNeedResponse {
	if b.ID.ID == "" {
		return response.BooksNeedResponse{}
	}

	year, _ := strconv.Atoi(b.Year)
	semseter, _ := strconv.Atoi(b.Semster)
	value, _ := strconv.ParseInt(b.Value, 10, 64)

	// Length year changes always at least equal to length supported years
	var yearChanges, supportedYears []int
	for i := 0; i < len(b.YearChanges); i++ {
		yearChange, _ := strconv.Atoi(b.YearChanges[i])
		yearChanges = append(yearChanges, yearChange)

		if i < len(b.SupportedYears) {
			supportedYear, _ := strconv.Atoi(b.SupportedYears[i])
			supportedYears = append(supportedYears, supportedYear)
		}
	}

	return response.BooksNeedResponse{
		ID:                  b.ID.ID,
		Year:                year,
		YearChanges:         yearChanges,
		Semster:             semseter,
		Value:               value,
		SupportedYears:      supportedYears,
		Donors:              b.Donors,
		ProvideDates:        b.ProvideDates,
		ProvidePeriods:      b.ProvidePeriods,
		ProvideStaffs:       b.ProvideStaffs,
		ProvideImageBlobIDs: b.ProvideImageBlobIDs,
		Donations:           b.Donations,
		WithdrawProposals:   b.WithdrawProposals,
		WithdrawsForNeed:    b.WithdrawsForNeed,
	}
}

func (h HealthInsuranceNeed) ToHealthInsuranceNeedReponse() response.HealthInsuranceNeedResponse {
	if h.ID.ID == "" {
		return response.HealthInsuranceNeedResponse{}
	}

	year, _ := strconv.Atoi(h.Year)
	value, _ := strconv.ParseInt(h.Value, 10, 64)

	// Length year changes always at least equal to length supported years
	var yearChanges, supportedYears []int
	for i := 0; i < len(h.YearChanges); i++ {
		yearChange, _ := strconv.Atoi(h.YearChanges[i])
		yearChanges = append(yearChanges, yearChange)

		if i < len(h.SupportedYears) {
			supportedYear, _ := strconv.Atoi(h.SupportedYears[i])
			supportedYears = append(supportedYears, supportedYear)
		}
	}

	return response.HealthInsuranceNeedResponse{
		ID:                  h.ID.ID,
		Year:                year,
		YearChanges:         yearChanges,
		Value:               value,
		SupportedYears:      supportedYears,
		Donors:              h.Donors,
		ProvideDates:        h.ProvideDates,
		ProvidePeriods:      h.ProvidePeriods,
		ProvideStaffs:       h.ProvideStaffs,
		ProvideImageBlobIDs: h.ProvideImageBlobIDs,
		Donations:           h.Donations,
		WithdrawProposals:   h.WithdrawProposals,
		WithdrawsForNeed:    h.WithdrawsForNeed,
	}
}
