package main

import (
	"bytes"
	"errors"
	"log"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/mshafiee/swephgo"
)

// House represents an astrological house
type House struct {
	SignName string
	Degree   float64
	Number   string
	DegreeUt float64
	Bodies   []int
}

var (
	lat, _    = strconv.ParseFloat(os.Getenv("LATITUDE"), 64)
	lon, _    = strconv.ParseFloat(os.Getenv("LONGITUDE"), 64)
	signNames = []string{"Aries", "Taurus", "Gemini", "Cancer", "Leo",
		"Virgo", "Libra", "Scorpio", "Sagittarius", "Capricorn", "Aquarius",
		"Pisces"}
	houseNames = []string{"0", "I", "II", "III", "IV", "V", "VI", "VII", "VIII",
		"IX", "X", "XI", "XII"}
	bodies = []int{
		swephgo.SeSun,
		swephgo.SeMoon,
		swephgo.SeMercury,
		swephgo.SeVenus,
		swephgo.SeMars,
		swephgo.SeJupiter,
		swephgo.SeSaturn,
		swephgo.SeUranus,
		swephgo.SeNeptune,
		swephgo.SePluto,
	}
	system = map[string]int{
		"Placidus":      int('P'),
		"Koch":          int('K'),
		"Porphyrius":    int('O'),
		"Regiomontanus": int('R'),
		"Equal":         int('E'),
		"Whole":         int('W'),
	}
)

/* Bodies() - return longitude of all planets
 */
func Bodies(when time.Time) []float64 {
	var b []float64
	for _, ipl := range bodies {
		x2, _ := Waldo(when, ipl, swephgo.SeflgSwieph+swephgo.SeflgRadians)
		b = append(b, x2[0])
	}
	return b
}

/* Houses() - fill in all houses (sign, position, cusp)
 */
func Houses(Cusps []float64) *[]House {
	var houses []House
	for house := 1; house <= 12; house++ {
		degreeUt := deg2rad(float64(Cusps[house]))
		for i, _ := range signNames {
			degLow := float64(i) * math.Pi / 6.0
			degHigh := float64((i + 1)) * math.Pi / 6.0
			if degreeUt >= degLow && degreeUt <= degHigh {
				houses = append(houses,
					House{
						SignName: signNames[i],
						Degree:   rad2deg(degreeUt - degLow),
						Number:   houseNames[house],
						DegreeUt: rad2deg(degreeUt),
					},
				)
			}
		}
	}
	return &houses
}

/* Cusps() gest cusps and asmc
 */
func Cusps(when time.Time, lat float64, lon float64, hsys int) ([]float64, []float64, error) {
	cusps := make([]float64, 13)
	asmc := make([]float64, 10)
	serr := make([]byte, 256)
	julianDay := julian(when)
	swephgo.SetTopo(lat, lon, 0)
	if eclflag := swephgo.Houses(*julianDay, lat, lon, hsys, cusps, asmc); eclflag == swephgo.Err {
		log.Printf("Error %d %s", eclflag, string(serr))
		return nil, nil, errors.New(string(serr))
	}
	return cusps, asmc, nil
}

/* Aspect() returns an aspect of two celectial bodies if any
or empty string
*/
func Aspect(body1 float64, body2 float64) string {
	aspect := ""
	angle := smallestSignedAngleBetween(body1, body2)
	if math.Abs(angle) < deg2rad(10.0) {
		aspect = "Conjunction"
	}
	if math.Abs(angle-math.Pi) < deg2rad(10.0) {
		aspect = "Opposition"
	}
	if math.Abs(angle-2.0*math.Pi/3.0) < deg2rad(8.0) {
		aspect = "Trine"
	}
	if math.Abs(angle-math.Pi/2.0) < deg2rad(6.0) {
		aspect = "Square"
	}
	if math.Abs(angle-math.Pi/3.0) < deg2rad(4.0) {
		aspect = "Sextile"
	}
	if math.Abs(angle-5.0*math.Pi/6.0) < deg2rad(2.0) {
		aspect = "Quincunx"
	}
	if math.Abs(angle-math.Pi/6.0) < deg2rad(1.0) {
		aspect = "Semi-sextile"
	}
	return aspect
}

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
	start = bod(start)
	end := start.AddDate(2, 0, 1) // look ahead up to 2 years and 1 day
	step := 0
	for d := start; d.After(end) == false; d = d.AddDate(0, 0, 1) {
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

/* SolarEclipse() find nearest solar eclipse
 */
func SolarEclipse(when time.Time, ifltype int) ([]float64, error) {
	julianDay := julian(when)
	x2 := make([]float64, 10)
	attr := make([]float64, 20)
	geopos := make([]float64, 10)
	serr := make([]byte, 256)
	method := swephgo.SeflgSwieph
	var eclflag int32
	tjdStart := *julianDay
	/* find next eclipse anywhere on Earth */
	if eclflag = swephgo.SolEclipseWhenGlob(tjdStart, method, ifltype, x2, 0, serr); eclflag == swephgo.Err {
		return x2, errors.New(string(serr))
	}
	/* the time of the greatest eclipse has been returned in tret[0];
	 * now we can find geographical position of the eclipse maximum */
	tjdStart = x2[0]
	if eclflag = swephgo.SolEclipseWhere(tjdStart, method, geopos, attr, serr); eclflag == swephgo.Err {
		return x2, errors.New(string(serr))
	}
	/* the geographical position of the eclipse maximum is in geopos[0] and geopos[1];
	 * now we can calculate the four contacts for this place. The start time is chosen
	 * a day before the maximum eclipse: */
	tjdStart = x2[0] - 1
	if eclflag = swephgo.SolEclipseWhenLoc(tjdStart, method, geopos, x2, attr, 0, serr); eclflag == swephgo.Err {
		return x2, errors.New(string(serr))
	}
	/* now x2[] contains the following values:
	 * x2[0] = time of greatest eclipse (Julian day number)
	 * x2[1] = first contact
	 * x2[2] = second contact
	 * x2[3] = third contact
	 * x2[4] = fourth contact */
	// Convert ecclipse back to Gregorian date
	return x2, nil
}

/*
search for any lunar eclipse, no matter which type
ifltype = 0;
search a total lunar eclipse
ifltype = SE_ECL_TOTAL;
search a partial lunar eclipse
ifltype = SE_ECL_PARTIAL;
search a penumbral lunar eclipse
ifltype = SE_ECL_PENUMBRAL;
*/
func LunarEclipse(when time.Time, eclType int) ([]float64, error) {
	julianDay := julian(when)
	// Fixed length array with results for eclipse calculation - so this is output
	x2 := make([]float64, 10)
	serr := make([]byte, 256)
	// Look for total eclipe for given julian date
	// method - 0 simple, 2 Swiss etc. look backward - No
	method := swephgo.SeflgSwieph
	backward := bool2int(false)
	if eclflag := swephgo.LunEclipseWhen(*julianDay, method, eclType, x2, backward, serr); eclflag == swephgo.Err {
		return x2, errors.New(string(serr))
	}
	return x2, nil
}

func getPlanetName(ipl int) string {
	pN := make([]byte, 15)
	swephgo.GetPlanetName(ipl, pN)
	pN = bytes.Trim(pN, "\x00") // to get rid of trailing NUL characters
	planetName := string(pN)
	return planetName
}

/* getHouse() get house for longitude in radians
given houses cusps
*/
func getHouse(rad float64, houses *[]House) string {
	for i := 0; i < len(*houses); i++ {
		degLow := deg2rad((*houses)[i].DegreeUt)
		var degHigh float64
		if i == len(*houses)-1 {
			degHigh = deg2rad((*houses)[0].DegreeUt)
		} else {
			degHigh = deg2rad((*houses)[i+1].DegreeUt)
		}
		if rad >= degLow && rad <= degHigh {
			return (*houses)[i].Number
		}
	}
	return (*houses)[0].Number
}

/* getSign() - cast longitude in radians to zodiac sign name
 */
func getSign(rad float64) string {
	for i, sign := range signNames {
		degLow := float64(i) * math.Pi / 6.0
		degHigh := float64((i + 1)) * math.Pi / 6.0
		if rad >= degLow && rad <= degHigh {
			return sign
		}
	}
	return ""
}
