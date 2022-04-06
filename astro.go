package main

import (
	"errors"
	"log"
	"math"
	"time"

	"github.com/mshafiee/swephgo"
)

/*
Wrapping up some Swiss Ephem code in Goland -
mind that guys at Astrodiesnt are old school C coding heros
*/
/*
What is the phase (ilumination) of a planet?
https://groups.io/g/swisseph/message/7327
*/
func Phase(when time.Time, planet int) (float64, error) {
	julianDay := julian(when)
	iflag := swephgo.SeflgSwieph // use SWISSEPH ephemeris, default
	attr := make([]float64, 20)
	serr := make([]byte, 256)
	if eclflag := swephgo.Pheno(*julianDay, planet, iflag, attr, serr); eclflag == swephgo.Err {
		log.Printf("Error %d %s", eclflag, string(serr))
		return 0.0, errors.New(string(serr))
	}
	return attr[1], nil
}

/*
Where is a planet (longitude, latitude, distance, speed in long., speed in lat., and speed in dist.)
*/
func Waldo(when time.Time, planet int, iflag int) ([]float64, error) {
	julianDay := julian(when)
	x2 := make([]float64, 6)
	serr := make([]byte, 256)
	if eclflag := swephgo.Calc(*julianDay, planet, iflag, x2, serr); eclflag == swephgo.Err {
		return x2, errors.New(string(serr))
	}
	return x2, nil
}

func RetroUt(start time.Time, ipl int, iflag int, jdx *float64, idir *int, serr *[]byte) int {
	var tx float64
	rval := Retro(start, ipl, iflag, &tx, idir, serr)
	if rval >= 0 {
		*jdx = tx - swephgo.Deltat(tx)
	}
	return rval
}

//int swe_next_direction_change(double jd0, int ipl, int iflag, double *jdx, int *idir, char *serr)
func Retro(start time.Time, ipl int, iflag int, jdx *float64, idir *int, serr *[]byte) int {
	// x2 := make([]float64, 6)
	var tx float64
	jd_step := 1.0
	jd0 := swephgo.Julday(start.Year(), int(start.Month()), start.Day(), float64(start.Hour()), swephgo.SeGregCal)
	x2, _ := Waldo(start, ipl, iflag)
	y0 := x2[0]
	y1 := x2[0]
	planetName := make([]byte, 10)
	swephgo.GetPlanetName(ipl, planetName)
	start = nod(start)
	end := start.AddDate(2, 0, 1) // look ahead up to 2 years and 1 day
	step := 0
	for d := start; d.After(end) == false; d = d.AddDate(0, 0, 1) {
		// fmt.Printf("%s %.10f retro date %s\n", string(planetName), y0, d.Format("2006-01-02 15:04:05"))
		jd := swephgo.Julday(d.Year(), int(d.Month()), d.Day(), float64(d.Hour()), swephgo.SeGregCal)
		x2, _ = Waldo(d, ipl, iflag)
		y2 := x2[0]
		// get parabola y = ax^2  + bx + c  and derivative y' = 2ax + b
		d1 := swephgo.Difdeg2n(y1, y0)
		d2 := swephgo.Difdeg2n(y2, y1)
		y0 = y1 // for next step
		y1 = y2
		b := (d1 + d2) / 2
		a := (d2 - d1) / 2
		if a == 0 {
			continue // curve is flat
		}
		tx = -b / a / 2.0 // time when derivative is zer0
		if tx < -1 || tx > 1 {
			continue
		}
		*jdx = jd - jd_step + tx*jd_step
		if *jdx-jd0 < 30.0/1440 {
			continue // ignore if within 30 minutes of start moment
		}
		// This is where magic happens
		for jd_step > 2/1440.0 {
			jd_step = jd_step / 2
			t1 := *jdx
			t0 := t1 - jd_step
			t2 := t1 + jd_step
			x2, _ = Waldo(jdToUTC(&t0), ipl, iflag)
			y0 = x2[0]
			x2, _ = Waldo(jdToUTC(&t1), ipl, iflag)
			y1 = x2[0]
			x2, _ = Waldo(jdToUTC(&t2), ipl, iflag)
			y2 = x2[0]
			d1 = swephgo.Difdeg2n(y1, y0)
			d2 = swephgo.Difdeg2n(y2, y1)
			b = (d1 + d2) / 2
			a = (d2 - d1) / 2
			if a == 0 {
				continue          // curve is flat }
				tx = -b / a / 2.0 // time when derivative is zer0
				if tx < -1 || tx > 1 {
					continue
				}
				*jdx = t1 + tx*jd_step
				tdiff := math.Abs(*jdx - t1)
				if tdiff < 1/86400.0 { // precision up to 1 minute
					break
				}
			}
			if a > 0 {
				*idir = 1
			} else {
				*idir = -1
			}
			step++
			return 0
		}
	}
	return 0
}

/* general helpers - should go to separete file */

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

func fixangle(a float64) float64 {
	return (a - 360*math.Floor(a/360))
}

func rad2deg(r float64) float64 {
	return (r * 180) / math.Pi
}

func deg2rad(d float64) float64 {
	return (d * math.Pi) / 180
}
