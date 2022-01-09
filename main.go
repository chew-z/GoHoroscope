package main

import (
	"fmt"

	"github.com/mshafiee/swephgo"
)

func main() {
	// Point to where Swiss Ephem files are located on your system
	// It is a good practice to do it as initialization
	// even when not using files
	swephgo.SetEphePath([]byte("/usr/local/share/sweph/ephe"))
	// Check version of library
	sweVer := make([]byte, 12)
	swephgo.Version(sweVer)
	fmt.Printf("Library used: Swiss Ephemeris v%s\n", sweVer)

	// Convert date from gregorian calendar to julian day (float)
	julianDay := swephgo.Julday(2021, 12, 31, 0, swephgo.SeGregCal)
	fmt.Printf("Julian day := %f\n", julianDay)

	// Turtles all the way down from here
	// swephgo is just baremetal, naked C
	// Use make to declare variables - preallocate space !!!

	// Fixed length array with results for eclipse calculation - so this is output
	tret := make([]float64, 10)
	// Placeholder for errors
	var serr []byte
	// Look for total eclipe for given julian date
	// method - 0 simple, 2 Swiss etc. look backward - No
	method := swephgo.SeflgSwieph
	backward := bool2int(false)
	result := swephgo.LunEclipseWhen(julianDay, method, swephgo.SeEclTotal, tret, backward, serr)
	fmt.Printf("Result is %d\n", result)
	// Print the result - lunar ecclipse date - tret[0] is time of maximum eclipse
	fmt.Println(tret[0])
	// Again pre-allocate space for outgoing values
	year := make([]int, 1)
	month := make([]int, 1)
	day := make([]int, 1)
	hour := make([]float64, 1)
	// Convert lunar ecclipse to back to Gregorian date
	swephgo.Revjul(tret[0], swephgo.SeGregCal, year, month, day, hour)
	// So when is the ecclipse?
	fmt.Printf("Eclipse date %d-%d-%d %.3f", year[0], month[0], day[0], hour[0])
}

func bool2int(b bool) int {
	if b {
		return 1
	}
	return 0
}
