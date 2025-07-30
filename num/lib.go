package num

import "strconv"

func ParseFloat64(numberStr string) float64 {
	num, err := strconv.ParseFloat(numberStr, 64)
	if err != nil {
		return 0
	}
	return num
}

func ParseInt64(numberStr string) int64 {
	num, err := strconv.ParseInt(numberStr, 10, 64)
	if err != nil {
		return 0
	}
	return num
}

func ParseInt(numberStr string) int {
	num, err := strconv.ParseInt(numberStr, 10, 32)
	if err != nil {
		return 0
	}
	return int(num)
}
