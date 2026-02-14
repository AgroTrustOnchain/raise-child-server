package response

import "time"

type MealNeedResponse struct {
	ID                   string              `json:"id"`
	Year                 int                 `json:"year"`
	Value                int64               `json:"value"`
	Donors               []string            `json:"donors"`
	Donations            []string            `json:"donations"`
	Durations            []WrapDuration      `json:"durations"`
	TotalSupportedMonths int                 `json:"total_supported_months"`
	SupportedYears       []WrapSupportedYear `json:"supported_years"`
	WithdrawProposals    []string            `json:"withdraw_proposals"`
	WithdrawsForNeed     []string            `json:"withdraws_for_need"`
}

type WrapDuration struct {
	StartPeriod time.Time `json:"start_period"`
	EndPeriod   time.Time `json:"end_period"`
}

type WrapSupportedYear struct {
	Year            int `json:"year"`
	SupportedMonths int `json:"supported_months"`
}
