package util

import "slices"

func StanderizeGender(gender string) string {
	var res string = StanderizeString(gender)

	var genders = []string{"nam", "nữ", "male", "female"}
	if existed := slices.Contains(genders, res); !existed {
		res = ""
	}

	return res
}
