package util

import "slices"

func StandardizeGender(gender string) string {
	var res string = StandardizeString(gender)

	var genders = []string{"nam", "nữ", "male", "female"}
	if existed := slices.Contains(genders, res); !existed {
		res = ""
	}

	return res
}
