package repository

import (
	"fmt"
	"math"
)

type generateRetrieveQueryRequest struct {
	table       string
	limitAmount int
	condition   string
	order       string
	page        int
	isGetCount  bool
}

// Caculate total pages
func caculateTotalPages(records, limitAmount int) int {
	return int(math.Ceil(float64(records) / float64(limitAmount)))
}

// Generate query to count the number of all records from a table from database
func generateCountTotalRecordsQuery(table, condition string) string {
	return "SELECT COUNT(*) FROM " + table + " WHERE " + condition
}

func generateRetrieveQuery(req generateRetrieveQueryRequest) string {
	if req.isGetCount {
		return generateCountTotalRecordsQuery(req.table, req.condition)
	}

	var res string = "SELECT * FROM " + req.table
	if req.condition != "" {
		res += " WHERE " + req.condition
	}

	if req.order != "" {
		res += req.order
	}

	var offSet int = (req.page - 1) * req.limitAmount

	return res + fmt.Sprintf(" LIMIT %d OFFSET %d", req.limitAmount, offSet)
}
