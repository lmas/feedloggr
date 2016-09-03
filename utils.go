package feedloggr2

import "time"

var today time.Time

func Now() time.Time {
	if today.IsZero() {
		today = time.Now()
	}
	return today
}
