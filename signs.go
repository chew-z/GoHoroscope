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

var (
	lat, _    = strconv.ParseFloat(os.Getenv("LATITUDE"), 64)
	lon, _    = strconv.ParseFloat(os.Getenv("LONGITUDE"), 64)
	numhouses = 12
	signNames = []string{"Aries", "Taurus", "Gemini", "Cancer", "Leo",
		"Virgo", "Libra", "Scorpio", "Sagittarius", "Capricorn", "Aquarius",
		"Pisces"}
	// Houses names
	hnames = []string{"0", "I", "II", "III", "IV", "V", "VI", "VII", "VIII",
		"IX", "X", "XI", "XII"}
	Cusps  = make([]float64, 13)
	Asmc   = make([]float64, 10)
	bodies = []int{
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
)

type aspectsetting struct {
	delta float64
	orb   float64
	title string
}

var aspectsettings = []aspectsetting{
	{180, 10, "Opposition"},
	{150, 2, "Quincunx"},
	{120, 8, "Trine"},
	{90, 6, "Square"},
	{60, 4, "Sextile"},
	{30, 1, "Semi-sextile"},
	{0, 10, "Conjunction"},
}

// House represents an astrological house cuspid
type House struct {
	SignName string
	Degree   float64
	Number   string
	DegreeUt float64
}

func init() {
	// Point to where Swiss Ephem files are located on your system
	// It is a good practice to do it as initialization
	// even when not using files
	swephgo.SetEphePath([]byte("/usr/local/share/sweph/ephe"))
}

func main() {
	now := time.Now().UTC()
	err := cusps(now, lat, lon)
	for i, c := range Cusps {
		fmt.Printf("House # %d - %.2f\n", i, c)
	}
	fmt.Println(Asmc[0])
	if err != nil {
		log.Println(err.Error())
	}
	// TODO - function
	H := Houses()
	for _, h := range *H {
		fmt.Printf("%s - %s - %.2f\n", h.Number, h.SignName, h.DegreeUt)
	}
	// TODO - function
	B := Bodies(now)
	for i, b1 := range B {
		fmt.Printf("%s - %.2f\n", getPlanetName(bodies[i]), rad2deg(b1))
		for j, b2 := range B[i+1:] {
			if a := aspect(b1, b2); a != "" {
				fmt.Printf("\t%s - %s - %.2f\n", a, getPlanetName(bodies[i+j+1]), rad2deg(b2))
			}
		}
	}

	swephgo.Close()
}

func Signs() {
	for x := 0.0; x < 2.0*math.Pi; x += math.Pi / 6.0 {
		fmt.Printf("Sign: %s\tbeg: %.3f, end: %.3f\tcosinus(beg): %.3f, cos(end): %.3f\n", sign(x), x, x+math.Pi/6.0, math.Cos(x), math.Cos(x+math.Pi/6.0))
	}

}

/* sign() - cast latitude in radians to zodiac sign name
 */
func sign(rad float64) string {
	for i, sign := range signNames {
		degLow := float64(i) * math.Pi / 6.0
		degHigh := float64((i + 1)) * math.Pi / 6.0
		if rad >= degLow && rad <= degHigh {
			return sign
		}
	}
	return ""
}

func Bodies(when time.Time) []float64 {
	var b []float64
	for _, ipl := range bodies {
		x2, _ := Waldo(when, ipl, swephgo.SeflgSwieph+swephgo.SeflgRadians)
		b = append(b, x2[0])
	}
	return b
}

func Houses() *[]House {
	var houses []House
	for house := 1; house <= numhouses; house++ {
		degreeUt := deg2rad(float64(Cusps[house]))
		for i, _ := range signNames {
			degLow := float64(i) * math.Pi / 6.0
			degHigh := float64((i + 1)) * math.Pi / 6.0
			if degreeUt >= degLow && degreeUt <= degHigh {
				houses = append(houses,
					House{
						SignName: signNames[i],
						Degree:   math.Round(rad2deg(degreeUt - degLow)),
						Number:   hnames[house],
						DegreeUt: math.Round(rad2deg(degreeUt)),
					},
				)
			}
		}
	}
	return &houses
}

func cusps(when time.Time, lat float64, lon float64) error {
	swephgo.SetTopo(lat, lon, 0)
	julianDay := julian(when)
	serr := make([]byte, 256)
	if eclflag := swephgo.Houses(*julianDay, lat, lon, int('P'), Cusps, Asmc); eclflag == swephgo.Err {
		log.Printf("Error %d %s", eclflag, string(serr))
		return errors.New(string(serr))
	}
	return nil
}

// makeAspect returns an Aspect for a given orb and two celectial bodies
func aspect(body1 float64, body2 float64) string {
	aspect := ""
	angle := smallestSignedAngleBetween(body1, body2)
	if math.Abs(angle) < deg2rad(10.0) {
		aspect = "Conjunction"
	}
	if math.Abs(angle-math.Pi) < deg2rad(10.0) {
		aspect = "Opposition"
	}
	if math.Abs(angle-math.Pi/2.0) < deg2rad(6.0) {
		aspect = "Square"
	}
	if math.Abs(angle-2.0*math.Pi/3.0) < deg2rad(8.0) {
		aspect = "Trine"
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

func smallestSignedAngleBetween(x float64, y float64) float64 {
	return math.Min(2.0*math.Pi-math.Abs(x-y), math.Abs(x-y))
}

// Make sure angle values are in within the 0 to 360 range
func normalize(angle float64) float64 {
	angle = math.Mod(angle, 2.0*math.Pi)
	if angle < 0 {
		angle += 2.0 * math.Pi
	}
	return angle
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
