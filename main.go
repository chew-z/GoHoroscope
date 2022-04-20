package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/araddon/dateparse"
	_ "github.com/joho/godotenv/autoload"
	"github.com/mshafiee/swephgo"
	"github.com/scylladb/termtables"
)

var (
	loc         string
	city        = os.Getenv("CITY")
	houseSystem = "Placidus"
	location    *time.Location
	swisspath   = os.Getenv("SWISSPATH")
)

func init() {
	houseSystem = os.Getenv("HOUSE_SYSTEM")
	location, _ = time.LoadLocation(city)
	swephgo.SetEphePath([]byte(swisspath))
}

func main() {
	now := time.Now()
	if len(os.Args) < 2 {
		fmt.Println("Usage horoscope -r|-e|-h [date]")
		return
	}
	if os.Args[1] == "-r" || os.Args[1] == "--retrograde" {
		PrintRetro(now, now.AddDate(1, 0, 1))
		return
	}
	if os.Args[1] == "-e" || os.Args[1] == "--eclipse" {
		PrintEclipse(now, now.AddDate(1, 0, 1))
		return
	}
	if os.Args[1] == "-h" || os.Args[1] == "--horoscope" {
		hsys := system[houseSystem]
		if len(os.Args) < 3 {
			PrintHoroscope(now, hsys) // lat, lon is given implicite in .env
		} else {
			when, err := dateparse.ParseLocal(os.Args[2])
			if err != nil {
				log.Println(err.Error())
				return
			}
			PrintHoroscope(when, hsys) // lat, lon is given implicite in .env
		}
	}
	defer swephgo.Close()
}

func PrintHoroscope(when time.Time, hsys int) {
	fmt.Printf("\n%s - lat: %.2f, lon: %.2f\n", when.In(location).Format(time.RFC822), lat, lon)
	if Cusps, Asmc, e := Cusps(when, lat, lon, hsys); e != nil {
		log.Panic(e)
	} else {
		fmt.Printf("Ascendant: %.2f MC: %.2f, House system: %s\n", Asmc[0], Asmc[1], houseSystem)
		fmt.Println()
		// TODO - function
		H := Houses(Cusps)
		table1 := termtables.CreateTable()
		table1.AddHeaders("House", "Position", "Cusp", "Sign")
		for _, h := range *H {
			// fmt.Printf("%s\t%.2f\t%.2f\t%s\n", h.Number, h.DegreeUt, h.Degree, h.SignName)
			table1.AddRow(h.Number, fmt.Sprintf("%.2f", h.DegreeUt), fmt.Sprintf("%.2f", h.Degree), h.SignName)
		}
		fmt.Println(table1.Render())
		// TODO - function
		B := Bodies(when)
		table2 := termtables.CreateTable()
		table2.AddHeaders("Planet", "Position", "House", "Sign", "Aspects")
		for i, b1 := range B {
			// fmt.Printf("House %s: %s - %.2f in %s\n", getHouse(b1, H), getPlanetName(bodies[i]), rad2deg(b1), getSign(b1))
			table2.AddRow(getPlanetName(bodies[i]), fmt.Sprintf("%.2f", rad2deg(b1)), getHouse(b1, H), getSign(b1))
			for j, b2 := range B[i+1:] {
				if asp := Aspect(b1, b2); asp != "" {
					c := fmt.Sprintf("%s %s in %s", asp, getPlanetName(bodies[i+j+1]), getSign(b2))
					table2.AddRow("", "", "", "", c)
				}
			}
		}
		fmt.Println(table2.Render())
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

	table := termtables.CreateTable()
	table.AddHeaders("Planet", "Starts", "Ends")
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
				d1 := prevDate.In(location).Format(time.RFC822)
				d2 := wd.In(location).Format(time.RFC822)
				// fmt.Printf("%s retrograde starts: %s ends: %s\n", planetName, d1, d2)
				table.AddRow(planetName, d1, d2)
			}
			prevDate = wd
			d = wd.AddDate(0, 0, 7) // start looking for next change in a direction 7 days ahead
		}
	}
	fmt.Println(table.Render())
}

func PrintEclipse(start time.Time, end time.Time) {
	d := start
	table1 := termtables.CreateTable()
	table1.AddHeaders("Lunar Eclipse")
	for d.After(end) == false {
		l, _ := LunarEclipse(d, swephgo.SeEclAlltypesLunar)
		wd := jdToUTC(&l[0])
		// fmt.Printf("Lunar eclipse: %s\t \n", jdToLocal(&l[0])) // eclipse maximum [0]
		table1.AddRow(jdToLocal(&l[0]))
		d = wd.AddDate(0, 0, 7)
	}
	d = start
	table2 := termtables.CreateTable()
	table2.AddHeaders("Solar Eclipse")
	for d.After(end) == false {
		s, _ := SolarEclipse(d, swephgo.SeEclAlltypesSolar)
		wd := jdToUTC(&s[0])
		// fmt.Printf("Solar eclipse: %s\n", jdToLocal(&s[0])) // eclipse maximum [0]
		table2.AddRow(jdToLocal(&s[0]))
		d = wd.AddDate(0, 0, 7)
	}
	fmt.Println(table1.Render())
	fmt.Println(table2.Render())
}
