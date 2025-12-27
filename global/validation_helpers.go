package global

import (
	"time"
)

func IsValidateTimeOfDay(timeOfDay string) bool {
	_, err := time.Parse("15:04", timeOfDay)
	if err != nil {
		return false
	} else if len(timeOfDay) < 5 {
		return false
	}

	return true
}
