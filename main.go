package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/mshafiee/swephgo"
)

var (
	// http.Clients should be reused instead of created as needed.
	client = &http.Client{
		Timeout: 5 * time.Second,
	}
	loc string
)

func init() {
	location, _ = time.LoadLocation(city)
	// Point to where Swiss Ephem files are located on your system
	// It is a good practice to do it as initialization
	// even when not using files
	swephgo.SetEphePath([]byte("/usr/local/share/sweph/ephe"))
}

func main() {
	// Draw some charts
	// http.HandleFunc("/", CloudCharts)
	// http.ListenAndServe(":8089", nil)

	// Test date conversion to julian and back
	// start := time.Now().UTC() // Start now
	// tx := julian(start)
	// rS := jdToLocal(tx)
	// rU := jdToUTC(tx)
	// fmt.Printf("%s\t\t%s\n", rS, rU)
	// Print table of retrograde movements
	PrintRetro()

	swephgo.Close()
}

/* PrintRetro - find retrograde movements
of planets for next 5 years
(when the movement is changing direction)
*/
func PrintRetro() {
	iflag := swephgo.SeflgSwieph
	var tx float64
	var idir int
	serr := make([]byte, 256)

	waldo := make([]float64, 6)
	var phase float64
	bodies := []int{
		swephgo.SeMercury,
		swephgo.SeVenus,
		swephgo.SeMars,
		swephgo.SeJupiter,
		swephgo.SeNeptune,
		swephgo.SeUranus,
		swephgo.SePluto,
	}
	start := time.Now().UTC()     // Start now
	end := start.AddDate(5, 0, 1) // and look ahead 5 years and 1 day
	for _, ipl := range bodies {
		planetName := make([]byte, 10)
		swephgo.GetPlanetName(ipl, planetName)
		d := start
		for d.After(end) == false {
			// find nearest change of direction
			if retval := RetroUt(d, ipl, iflag, &tx, &idir, &serr); retval < 0 {
				log.Printf("Error %s", string(serr))
				return
			}
			// what is the vector?
			direction := "direct"
			if idir < 0 {
				direction = "retro"
			}
			wd := jdToUTC(&tx)
			waldo, _ = Waldo(wd, ipl, swephgo.SeflgSwieph)
			phase, _ = Phase(wd, ipl)
			fmt.Printf("%s\t%s\t%s\t-\t%.5f\t%.1f\t%.5f\n", jdToLocal(&tx), string(planetName), direction, phase, waldo[0], waldo[2])
			d = wd.AddDate(0, 0, 7) // start looking for next change in a direction 7 days ahead
		}
	}
}
