package main

import (
	"math"
	"time"

	"github.com/mshafiee/swephgo"
)

/* general helpers */

func jdToUTC(jd *float64) time.Time {
	year := make([]int, 1)
	month := make([]int, 1)
	day := make([]int, 1)
	hour := make([]float64, 1)
	swephgo.Revjul(*jd, swephgo.SeGregCal, year, month, day, hour)
	h := int(hour[0])
	m := int(60 * (hour[0] - float64(h)))
	utc := time.Date(year[0], time.Month(month[0]), day[0], h, m, 0, 0, time.UTC)
	return utc
}

func jdToLocal(jd *float64) time.Time {
	utc := jdToUTC(jd)
	return utc.In(location)
}

func julian(d time.Time) *float64 {
	h := float64(d.Hour()) + float64(d.Minute())/60 + float64(d.Second())/3600
	jd := swephgo.Julday(d.Year(), int(d.Month()), d.Day(), h, swephgo.SeGregCal)
	return &jd
}

/* Begining of the Day */
func bod(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.Local)
}

/* Noon of the Day */
func nod(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 12, 0, 0, 0, time.Local)
}

func smallestSignedAngleBetween(x float64, y float64) float64 {
	return math.Min(2.0*math.Pi-math.Abs(x-y), math.Abs(x-y))
}

func fixangle(a float64) float64 {
	return (a - 360*math.Floor(a/360))
}

func rad2deg(r float64) float64 {
	return (r * 180) / math.Pi
}

func deg2rad(d float64) float64 {
	return (d * math.Pi) / 180
}

func bool2int(b bool) int {
	if b {
		return 1
	}
	return 0
}
