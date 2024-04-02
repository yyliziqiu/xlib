package xtime

import (
	"fmt"
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	loc := time.Local
	tt := time.Date(2024, 4, 2, 19, 27, 31, 0, loc)

	dt0, dt1 := DayRange(tt, loc)
	fmt.Println(dt0.Format(time.DateTime), dt1.Format(time.DateTime))

	wt0, wt1 := WeekRange(tt, loc)
	fmt.Println(wt0.Format(time.DateTime), wt1.Format(time.DateTime))

	mt0, mt1 := MonthRange(tt, loc)
	fmt.Println(mt0.Format(time.DateTime), mt1.Format(time.DateTime))

	yt0, yt1 := YearRange(tt, loc)
	fmt.Println(yt0.Format(time.DateTime), yt1.Format(time.DateTime))
}
