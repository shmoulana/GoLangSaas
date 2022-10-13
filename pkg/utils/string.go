package utils

import "strconv"

func StringToInt(val string, def int) int {
	var result int

	result, err := strconv.Atoi(val)
	if err != nil {
		return def
	}

	return result
}
