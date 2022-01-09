package main

import (
	"fmt"

	"github.com/mshafiee/swephgo"
)

func main() {
	swephgo.SetEphePath([]byte("/usr/local/share/sweph/ephe"))

	sweVer := make([]byte, 12)
	swephgo.Version(sweVer)
	fmt.Printf("Library used: Swiss Ephemeris v%s\n", sweVer)

	julianDay := swephgo.Julday(2021, 12, 31, 0, swephgo.SeGregCal)
	fmt.Printf("Julian day := %f\n", julianDay)
	tret := make([]float64, 10)
	var serr []byte
	result := swephgo.LunEclipseWhen(julianDay, 0, swephgo.SeEclTotal, tret, 0, serr)
	fmt.Printf("Result is %d\n", result)
	fmt.Println(tret[0])
	year := make([]int, 1)
	month := make([]int, 1)
	day := make([]int, 1)
	hour := make([]float64, 1)
	swephgo.Revjul(tret[0], swephgo.SeGregCal, year, month, day, hour)

	fmt.Printf("Eclipse date %d-%d-%d %.3f", year[0], month[0], day[0], hour[0])
}
