package entities

type EditNeedDates struct {
	ID        ID     `json:"id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}
