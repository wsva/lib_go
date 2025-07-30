package math

func AbsInt(i int) int {
	if i < 0 {
		return i * -1
	}
	return i
}

func AbsInt64(i int64) int64 {
	if i < 0 {
		return i * -1
	}
	return i
}
