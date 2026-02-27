package util

func StanderizeSortOrder(order string) string {
	var res string
	order = StanderizeString(order)
	switch order {
	case "desc":
		res = "DESC"
	case "asc":
		res = "ASC"
	}

	return res
}
