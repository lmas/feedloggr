package feedloggr2

import "time"

var Today time.Time

func Now() time.Time {
	if Today.IsZero() {
		Today = time.Now()
	}
	return Today
}
