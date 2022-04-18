package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"
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

func init() {
	swephgo.SetEphePath([]byte("/usr/local/share/sweph/ephe"))
}

func main() {
	now := time.Now()
	hsys := system["Placidus"]
	PrintHoroscope(now, hsys) // lat, lon is given implicite in .env
	defer swephgo.Close()
}

func PrintHoroscope(when time.Time, hsys int) {
	fmt.Printf("%s - lat: %.2f, lon: %.2f\n", when.Format(time.RFC822), lat, lon)
	if Cusps, Asmc, e := Cusps(when, lat, lon, hsys); e != nil {
		log.Panic(e)
	} else {
		fmt.Printf("Ascendant: %.2f MC: %.2f\n", Asmc[0], Asmc[1])
		// TODO - function
		H := Houses(Cusps)
		for _, h := range *H {
			fmt.Printf("%s\t%.2f\t%.2f\t%s\n", h.Number, h.DegreeUt, h.Degree, h.SignName)
		}
		fmt.Println()
		// TODO - function
		B := Bodies(when)
		for i, b1 := range B {
			fmt.Printf("House %s: %s - %.2f in %s\n", getHouse(b1, H), getPlanetName(bodies[i]), rad2deg(b1), getSign(b1))
			for j, b2 := range B[i+1:] {
				if asp := Aspect(b1, b2); asp != "" {
					fmt.Printf("\t%s - %s - %.2f in %s\n", asp, getPlanetName(bodies[i+j+1]), rad2deg(b2), getSign(b2))
				}
			}
		}
	}
}

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

func Waldo(when time.Time, planet int, iflag int) ([]float64, error) {
	julianDay := julian(when)
	x2 := make([]float64, 6)
	serr := make([]byte, 256)
	if eclflag := swephgo.Calc(*julianDay, planet, iflag, x2, serr); eclflag == swephgo.Err {
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

func smallestSignedAngleBetween(x float64, y float64) float64 {
	return math.Min(2.0*math.Pi-math.Abs(x-y), math.Abs(x-y))
}

func julian(d time.Time) *float64 {
	h := float64(d.Hour()) + float64(d.Minute())/60 + float64(d.Second())/3600
	jd := swephgo.Julday(d.Year(), int(d.Month()), d.Day(), h, swephgo.SeGregCal)
	return &jd
}

func rad2deg(r float64) float64 {
	return (r * 180) / math.Pi
}

func deg2rad(d float64) float64 {
	return (d * math.Pi) / 180
}
