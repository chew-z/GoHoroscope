package main

import (
	"fmt"
	"log"
	"math"
	"time"

	"github.com/mshafiee/swephgo"
)

type Date struct {
	Year   int
	Month  int
	Day    int
	Hour   int
	Minute int
}

func init() {
	swephgo.SetEphePath([]byte("/usr/local/share/sweph/ephe"))
}

/*
Adapted from example code
https://s3-us-west-1.amazonaws.com/groupsioattachments/3024/64370695/7781/0?AWSAccessKeyId=AKIAJECNKOVMCCU3ATNQ&Expires=1649102211&Signature=n0Mde4ILrkzNVr%2BJFUOilkdh%2FWI%3D&response-content-disposition=inline%3B+filename%3D%22st12.c%22
*/
func main() {
	iflag := swephgo.SeflgSwieph
	var tx float64
	var idir int
	serr := make([]byte, 256)
	bodies := []int{
		swephgo.SeMercury,
		swephgo.SeVenus,
		swephgo.SeMars,
		swephgo.SeJupiter,
		swephgo.SeNeptune,
		swephgo.SeUranus,
		swephgo.SePluto,
	}
	for _, ipl := range bodies {
		planetName := make([]byte, 10)
		swephgo.GetPlanetName(ipl, planetName)
		start := time.Now().UTC()     // Start now
		end := start.AddDate(5, 0, 1) // and look ahead 5 years and 1 day
		d := start
		for d.After(end) == false {
			// find nearest change of direction
			retval := RetroUt(d, ipl, iflag, &tx, &idir, &serr)
			if retval < 0 {
				log.Printf("Error %s", string(serr))
				return
			}
			// what is the vector
			direction := "direct"
			if idir < 0 {
				direction = "retro"
			}
			rD := christian(tx)
			rT := chrisToLocal(&rD)
			fmt.Printf("%s\t%s\t%s\n", string(planetName), direction, rT)
			d = chrisToUTC(&rD).AddDate(0, 0, 7) // start looking for next change in a direction 7 days ahead
		}
	}
	swephgo.Close()
}

// jdx   must be pointer to double, returns moment of direction change
// idir  must be pointer to integer, returns -1 if objects gets retrograde, 1 if direct
// The start moment jd0_ut must be at least 30 minutes before the direction change.
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
	jd_step := 1.0
	x2 := make([]float64, 6)
	var y1, tx float64
	if jd_step <= 0 {
		jd_step = 1.0
	}
	jd0 := swephgo.Julday(start.Year(), int(start.Month()), start.Day(), float64(start.Hour()), swephgo.SeGregCal)
	rval := swephgo.Calc(jd0, ipl, iflag, x2, *serr)
	if rval < 0 {
		return int(rval)
	}
	y0 := x2[0]
	y1 = x2[0]
	planetName := make([]byte, 10)
	start = nod(start)
	end := start.AddDate(2, 0, 1) // look ahead up to 2 years and 1 day
	step := 0
	for d := start; d.After(end) == false; d = d.AddDate(0, 0, 1) {
		swephgo.GetPlanetName(ipl, planetName)
		// fmt.Printf("%s %.10f retro date %s\n", string(planetName), y0, d.Format("2006-01-02 15:04:05"))
		jd := swephgo.Julday(d.Year(), int(d.Month()), d.Day(), float64(d.Hour()), swephgo.SeGregCal)
		rval = swephgo.Calc(jd, ipl, iflag, x2, *serr)
		if rval < 0 {
			return int(rval)
		}
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
			rval = swephgo.Calc(t0, ipl, iflag, x2, *serr)
			if rval < 0 {
				return int(rval)
			}
			y0 = x2[0]
			rval = swephgo.Calc(t1, ipl, iflag, x2, *serr)
			if rval < 0 {
				return int(rval)
			}
			y1 = x2[0]
			rval = swephgo.Calc(t2, ipl, iflag, x2, *serr)
			if rval < 0 {
				return int(rval)
			}
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
			return int(rval)
		}
	}
	return 0
}

func christian(tret float64) Date {
	var dt Date
	year := make([]int, 1)
	month := make([]int, 1)
	day := make([]int, 1)
	hour := make([]float64, 1)
	// Convert back to Gregorian date
	swephgo.Revjul(tret, swephgo.SeGregCal, year, month, day, hour)
	dt.Year = year[0]
	dt.Month = month[0]
	dt.Day = day[0]
	h := int(hour[0])
	dt.Hour = h
	m := int(60 * (hour[0] - float64(h)))
	dt.Minute = m
	return dt
}
func chrisToUTC(rD *Date) time.Time {
	utc := time.Date(rD.Year, time.Month(rD.Month), rD.Day, rD.Hour, rD.Minute, 0, 0, time.UTC)
	return utc
}

func chrisToLocal(rD *Date) time.Time {
	local := time.Date(rD.Year, time.Month(rD.Month), rD.Day, rD.Hour, rD.Minute, 0, 0, time.Local)
	return local
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
