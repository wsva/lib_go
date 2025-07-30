package str

import (
	"fmt"
	"strings"
)

// GetSubString comment
func GetSubString(str string, start, end int) (string, error) {
	rs := []rune(str)
	length := len(rs)

	if start < 0 || start > length {
		return "", fmt.Errorf("start is wrong")
	}

	if end < start || end > length {
		return "", fmt.Errorf("end is wrong")
	}

	return string(rs[start:end]), nil
}

func CompareStringMap(m1, m2 map[string]string) bool {
	for k, v := range m1 {
		if _, found := m2[k]; !found {
			return false
		}
		if m2[k] != v {
			return false
		}
		delete(m2, k)
	}
	for k, v := range m2 {
		if _, found := m1[k]; !found {
			return false
		}
		if m1[k] != v {
			return false
		}
	}
	return true
}

func StringListHas(list []string, item string) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}
	return false
}

var TemplateExecuteReplacer = strings.NewReplacer(
	"&", "&amp;",
	"<", "&lt;",
	">", "&gt;",
	// "&#34;" is shorter than "&quot;".
	`"`, "&#34;",
	// "&#39;" is shorter than "&apos;" and apos was not in HTML until HTML5.
	"'", "&#39;",
)
