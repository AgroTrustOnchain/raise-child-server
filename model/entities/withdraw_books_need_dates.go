package entities

type BooksNeedWithdrawDates struct {
	ID                 ID     `json:"id"`
	FirstSemesterDate  string `json:"first_semester_date"`
	SecondSemesterDate string `json:"second_semester_date"`
}

type HealthInsuranceNeedWithdrawDate struct {
	ID           ID     `json:"id"`
	ExpectedDate string `json:"expected_date"`
}
