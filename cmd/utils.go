package main

import (
	"fmt"
	"math"
	"time"
)

func afterDur(then int64) string {
	if then == 0 {
		return "-"
	}
	now := time.Now().Unix()
	diff := then - now
	units := []string{"year", "month", "day", "hour", "minute", "second"}
	unitSecs := []int64{31536000, 2592000, 86400, 3600, 60, 1}
	for i, unit := range units {
		if diff > unitSecs[i] || unit == "seconds" {
			mults := math.Round(float64(diff) / float64(unitSecs[i]))
			if mults > 1 {
				unit += "s"
			}
			return fmt.Sprintf("in %v %v", mults, unit)
		}
	}
	return "Something went wrong!"
}
