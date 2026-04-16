package util

import "slices"

func StandardizeRelation(relation string) string {
	var res string = StandardizeString(relation)

	var relations = []string{"father", "mother", "guardian"}
	if existed := slices.Contains(relations, res); !existed {
		res = ""
	}

	return res
}
