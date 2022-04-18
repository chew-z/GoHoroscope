package main

import (
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/mshafiee/swephgo"
)

var (
	loc         string
	city        = os.Getenv("CITY")
	location, _ = time.LoadLocation(city)
)

func init() {
	swephgo.SetEphePath([]byte("/usr/local/share/sweph/ephe"))
}

func main() {
	now := time.Now()
	PrintRetro(now, now.AddDate(1, 0, 1))
	PrintEclipse(now, now.AddDate(1, 0, 1))
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

/* PrintRetro - find retrograde movements
(when the movement is changing direction)
*/
func PrintRetro(start time.Time, end time.Time) {
	iflag := swephgo.SeflgSwieph
	var tx float64
	var idir int
	serr := make([]byte, 256)

	for _, ipl := range bodies {
		if ipl < 2 {
			continue
		}
		planetName := getPlanetName(ipl)
		d := start
		i := 0
		prevDate := start
		for d.After(end) == false {
			// find nearest change of direction
			if retval := RetroUt(d, ipl, iflag, &tx, &idir, &serr); retval < 0 {
				log.Printf("Error %s", string(serr))
				return
			}
			wd := jdToUTC(&tx)
			if i == 0 { // skip first step
				i++
				continue
			}
			// what is the vector?
			if idir > 0 {
				fmt.Printf("%s retrograde starts: %s ends: %s\n", planetName, prevDate.Format(time.RFC822), wd.Format(time.RFC822))
			}
			prevDate = wd
			d = wd.AddDate(0, 0, 7) // start looking for next change in a direction 7 days ahead
		}
		fmt.Println()
	}
}

func PrintEclipse(start time.Time, end time.Time) {
	d := start
	for d.After(end) == false {
		l, _ := LunarEclipse(d, swephgo.SeEclAlltypesLunar)
		wd := jdToUTC(&l[0])
		fmt.Printf("Lunar eclipse: %s\t \n", jdToLocal(&l[0])) // eclipse maximum [0]
		d = wd.AddDate(0, 0, 7)
	}
	fmt.Println()
	d = start
	for d.After(end) == false {
		s, _ := SolarEclipse(d, swephgo.SeEclAlltypesSolar)
		wd := jdToUTC(&s[0])
		fmt.Printf("Solar eclipse: %s\n", jdToLocal(&s[0])) // eclipse maximum [0]
		d = wd.AddDate(0, 0, 7)
	}
	fmt.Println()
}
