package feedloggr2

import "time"

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func inList(s string, l []string) bool {
	for _, ss := range l {
		if ss == s {
			return true
		}
	}
	return false
}

func date(t time.Time) string {
	return t.Format("2006-01-02")
}
