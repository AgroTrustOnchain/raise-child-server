package util

func StandardizeSortOrder(order string) string {
	var res string
	order = StandardizeString(order)
	switch order {
	case "desc":
		res = "DESC"
	case "asc":
		res = "ASC"
	default:
		res = "DESC"
	}

	return res
}

func StandardizeSortCriteria(sc string) string {
	var res string
	sc = StandardizeString(sc)
	switch sc {
	case "created_at":
		res = sc
	case "target", "withdraw_amount", "amount":
		res = sc
	default:
		sc = "created_at"
	}

	return res
}
