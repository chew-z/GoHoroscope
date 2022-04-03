package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/mshafiee/swephgo"
)

/*Date - ..
 */
type Date struct {
	Year   int
	Month  int
	Day    int
	Hour   int
	Minute int
}

var (
	birthDay, birthMonth, birthYear int
	birthHour, latitude, longitude  float64
	location                        string
)

func init() {
	birthDay, _ = strconv.Atoi(os.Getenv("BIRTHDAY"))
	birthMonth, _ = strconv.Atoi(os.Getenv("BIRTHMONTH"))
	birthYear, _ = strconv.Atoi(os.Getenv("BIRTHYEAR"))
	birthHour, _ = strconv.ParseFloat(os.Getenv("BIRTHHOUR"), 64)
	latitude, _ = strconv.ParseFloat(os.Getenv("LATITUDE"), 64)
	longitude, _ = strconv.ParseFloat(os.Getenv("LONGITUDE"), 64)
	location = os.Getenv("LOCATION")
}

func main() {
	// Point to where Swiss Ephem files are located on your system
	// It is a good practice to do it as initialization
	// even when not using files
	swephgo.SetEphePath([]byte("/usr/local/share/sweph/ephe"))
	loc, _ := time.LoadLocation(location)
	start := time.Now().UTC()
	start = Bod(start)
	end := start.AddDate(0, 1, 0)
	bodies := []int{swephgo.SeSun, swephgo.SeMoon, swephgo.SeMercury, swephgo.SeVenus}
	for i, ipl := range bodies {
		planetName := make([]byte, 10)
		swephgo.GetPlanetName(i, planetName)
		fmt.Printf("---- %s Phenomen---\n", string(planetName))
		fmt.Printf("date\tphase\tlongitude\tlatitude\n")
		for d := start; d.After(end) == false; d = d.AddDate(0, 0, 1) {
			julianDay := swephgo.Julday(d.Year(), int(d.Month()), d.Day(), float64(d.Hour()), swephgo.SeGregCal)
			fmt.Printf(d.In(loc).Format("2006-01-02 15:04 "))
			p, _ := phase(&julianDay, ipl)
			ll, _ := waldo(&julianDay, ipl)
			fmt.Printf("%.3f\t %.3f\t%.3f\n", p, ll[0], ll[1])
		}
	}
	swephgo.Close()
}

/*
https://groups.io/g/swisseph/message/7327
*/
func phase(julianDay *float64, planet int) (float64, error) {
	iflag := swephgo.SeflgSwieph // use SWISSEPH ephemeris, default
	attr := make([]float64, 20)
	serr := make([]byte, 256)
	eclflag := swephgo.Pheno(*julianDay, planet, iflag, attr, serr)
	if eclflag == swephgo.Err {
		log.Printf("Error %d %s", eclflag, string(serr))
		return 0.0, errors.New(string(serr))
	}
	return attr[1], nil
}

/* Where is a planet (lon, lat)

 */
func waldo(julianDay *float64, planet int) ([]float64, error) {
	lonlat := make([]float64, 2)
	iflag := swephgo.SeflgSwieph // use SWISSEPH ephemeris, default
	x2 := make([]float64, 6)
	serr := make([]byte, 256)
	eclflag := swephgo.Calc(*julianDay, planet, iflag, x2, serr)
	if eclflag == swephgo.Err {
		return lonlat, errors.New(string(serr))
	}
	lonlat[0] = x2[0]
	lonlat[1] = x2[1]
	return lonlat, nil
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

func Bod(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

func bool2int(b bool) int {
	if b {
		return 1
	}
	return 0
}
