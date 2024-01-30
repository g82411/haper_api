package date

import (
	"time"
)

// GetNowDateString Timezome: Asia/Taipei
func GetNowDateString() string {
	loc, err := time.LoadLocation("Asia/Taipei")
	if err != nil {
		panic(err)
	}
	now := time.Now().In(loc)
	dateStr := now.Format("2006-01-02")
	return dateStr
}
