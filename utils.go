package feedloggr2

import "time"

const Version string = "0.1"

var today time.Time

func Now() time.Time {
	if today.IsZero() {
		today = time.Now()
	}
	return today
}
