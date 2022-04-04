package main

import (
	"errors"
	"log"
	"net/http"
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
	// http.Clients should be reused instead of created as needed.
	client = &http.Client{
		Timeout: 5 * time.Second,
	}
)

func init() {
	// Point to where Swiss Ephem files are located on your system
	// It is a good practice to do it as initialization
	// even when not using files
	swephgo.SetEphePath([]byte("/usr/local/share/sweph/ephe"))
}

func main() {
	// Where the magic happens
	http.HandleFunc("/", CloudCharts)
	http.ListenAndServe(":8089", nil)

	swephgo.Close()
}

/*
What is a phase (ilumination) of a planet?
https://groups.io/g/swisseph/message/7327
*/
func Phase(julianDay *float64, planet int) (float64, error) {
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
