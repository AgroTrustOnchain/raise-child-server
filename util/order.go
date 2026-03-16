package util

func StanderizeSortOrder(order string) string {
	var res string
	order = StanderizeString(order)
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

func StanderizeSortCriteria(sc string) string {
	var res string
	sc = StanderizeString(sc)
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
