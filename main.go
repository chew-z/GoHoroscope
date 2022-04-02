package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

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
)

func init() {
	birthDay, _ = strconv.Atoi(os.Getenv("BIRTHDAY"))
	birthMonth, _ = strconv.Atoi(os.Getenv("BIRTHMONTH"))
	birthYear, _ = strconv.Atoi(os.Getenv("BIRTHYEAR"))
	birthHour, _ = strconv.ParseFloat(os.Getenv("BIRTHHOUR"), 64)
	latitude, _ = strconv.ParseFloat(os.Getenv("LATITUDE"), 64)
	longitude, _ = strconv.ParseFloat(os.Getenv("LONGITUDE"), 64)

}

func main() {
	// Point to where Swiss Ephem files are located on your system
	// It is a good practice to do it as initialization
	// even when not using files
	swephgo.SetEphePath([]byte("/usr/local/share/sweph/ephe"))
	// Check version of library
	sweVer := make([]byte, 12)
	swephgo.Version(sweVer)
	fmt.Printf("Library used: Swiss Ephemeris v%s\n", sweVer)

	var julianDay float64
	// cT := time.Now().UTC()
	// Convert date from gregorian calendar to julian day (float)
	// julianDay := swephgo.Julday(cT.Year(), int(cT.Month()), cT.Day(), float64(cT.Hour()), swephgo.SeGregCal)
	julianDay = swephgo.Julday(birthYear, birthMonth, birthDay, birthHour, swephgo.SeGregCal)
	fmt.Printf("Julian day := %f\n", julianDay)

	// Turtles all the way down from here
	// swephgo is just baremetal, naked C
	// Use make to declare variables - preallocate space !!!

	planets(&julianDay)

	ifltype := swephgo.SeEclTotal
	lunarEclipse(&julianDay, ifltype)

	// ifltype := swephgo.SeEclTotal
	ifltype = swephgo.SeEclAlltypesSolar
	solarEclipse(&julianDay, ifltype)

	julianDay = swephgo.Julday(birthYear, birthMonth, birthDay, birthHour, swephgo.SeGregCal)
	// Convert ecclipse back to Gregorian date
	birthdate := christian(julianDay)
	fmt.Printf("Birth date %d-%d-%d %2d:%2d\n", birthdate.Year, birthdate.Month, birthdate.Day, birthdate.Hour, birthdate.Minute)

	houses(&julianDay)

	swephgo.Close()
}

// func housePos() {

// 	hsys := int('W')
// 	xpin := make([]float64, 2)
// 	serr := make([]byte, 256)

// 	swephgo.HousePos(hsys, xpin, serr)
// }

func houses(julianDay *float64) {
	td := *julianDay
	cusps := make([]float64, 13)
	ascmc := make([]float64, 10)
	serr := make([]byte, 256)
	eclflag := swephgo.Houses(td, latitude, longitude, int('P'), cusps, ascmc)

	if eclflag == swephgo.Err {
		log.Printf("Error %d %s", eclflag, string(serr))
	}

	fmt.Println("---- Houses ---")
	for i, c := range cusps {
		fmt.Println(i, c)
	}
	fmt.Println()
	for i, a := range ascmc {
		fmt.Println(i, a)
	}

}

func solarEclipse(julianDay *float64, ifltype int) {
	var eclipse Date
	tret := make([]float64, 10)
	attr := make([]float64, 20)
	geopos := make([]float64, 10)
	// Placeholder for errors
	serr := make([]byte, 256)
	/* 0 default ephemeris, 2 - Swiss */
	method := swephgo.SeflgSwieph
	var eclflag int32

	fmt.Println("---- Nearest Solar eclipse ---")
	tjdStart := *julianDay
	/* find next eclipse anywhere on Earth */
	eclflag = swephgo.SolEclipseWhenGlob(tjdStart, method, ifltype, tret, 0, serr)
	if eclflag == swephgo.Err {
		log.Printf("Error %d %s", eclflag, string(serr))
	}
	/* the time of the greatest eclipse has been returned in tret[0];
	 * now we can find geographical position of the eclipse maximum */
	tjdStart = tret[0]
	eclflag = swephgo.SolEclipseWhere(tjdStart, method, geopos, attr, serr)
	if eclflag == swephgo.Err {
		log.Printf("Error %d %s", eclflag, string(serr))
	}
	/* the geographical position of the eclipse maximum is in geopos[0] and geopos[1];
	 * now we can calculate the four contacts for this place. The start time is chosen
	 * a day before the maximum eclipse: */
	tjdStart = tret[0] - 1
	eclflag = swephgo.SolEclipseWhenLoc(tjdStart, method, geopos, tret, attr, 0, serr)
	if eclflag == swephgo.Err {
		log.Printf("Error %d %s", eclflag, string(serr))
	}
	/* now tret[] contains the following values:
	 * tret[0] = time of greatest eclipse (Julian day number)
	 * tret[1] = first contact
	 * tret[2] = second contact
	 * tret[3] = third contact
	 * tret[4] = fourth contact */
	// Convert ecclipse back to Gregorian date
	eclipse = christian(tret[0])
	fmt.Printf("Solar eclipse %d-%d-%d %d:%d\n", eclipse.Year, eclipse.Month, eclipse.Day, eclipse.Hour, eclipse.Minute)
}

/*
double tjd_start = start date for search, Jul. day UT
int32 ifl  = ephemeris flag
int32 ifltype = eclipse type wanted: SE_ECL_TOTAL etc. or 0, if any eclipse type
double *tret = return array, 10 doubles, see below
AS_BOOL backward = TRUE, if backward search
char *serr = return error string
*/
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
func lunarEclipse(julianDay *float64, eclType int) {

	fmt.Println("---- Nearest Lunar eclipse ---")
	// Fixed length array with results for eclipse calculation - so this is output
	tret := make([]float64, 10)
	// Placeholder for errors
	serr := make([]byte, 256)
	// Look for total eclipe for given julian date
	// method - 0 simple, 2 Swiss etc. look backward - No
	method := swephgo.SeflgSwieph
	backward := bool2int(false)
	eclflag := swephgo.LunEclipseWhen(*julianDay, method, eclType, tret, backward, serr)
	if eclflag == swephgo.Err {
		log.Printf("Error %d %s", eclflag, string(serr))
	}
	// lunar ecclipse date - tret[0] is time of maximum eclipse
	// Convert lunar ecclipse back to Gregorian date
	var eclipse Date
	eclipse = christian(tret[0])
	// So when is the ecclipse?
	fmt.Printf("Eclipse maximum %d-%d-%d %d:%d\n", eclipse.Year, eclipse.Month, eclipse.Day, eclipse.Hour, eclipse.Minute)
	eclipse = christian(tret[4])
	fmt.Printf("Eclipse totality begins %d-%d-%d %d:%d\n", eclipse.Year, eclipse.Month, eclipse.Day, eclipse.Hour, eclipse.Minute)
	eclipse = christian(tret[5])
	fmt.Printf("Eclipse totality ends %d-%d-%d %d:%d\n", eclipse.Year, eclipse.Month, eclipse.Day, eclipse.Hour, eclipse.Minute)

}

/*
tjd_ut   = Julian day, Universal Time
ipl      = body number
iflag    = a 32 bit integer containing bit flags that indicate what kind of computation is wanted
xx       = array of 6 doubles for longitude, latitude, distance, speed in long., speed in lat., and speed in dist.
serr[256] = character string to return error messages in case of error. */
func planets(julianDay *float64) {
	iflag := swephgo.SeflgSwieph // use SWISSEPH ephemeris, default
	fmt.Println("---- List planets ---")
	fmt.Println("planet - longitude, latitude, distance")
	for i := swephgo.SeSun; i <= swephgo.SeVesta; i++ {
		planet := make([]byte, 20)
		x2 := make([]float64, 6)
		serr := make([]byte, 256)
		if i == swephgo.SeEarth {
			continue
		}
		swephgo.GetPlanetName(i, planet)
		swephgo.Calc(*julianDay, i, iflag, x2, serr)
		// fmt.Println(serr)
		fmt.Printf("%s - %.3f %.3f %.3f\n", string(planet), x2[0], x2[1], x2[2])
	}
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

func bool2int(b bool) int {
	if b {
		return 1
	}
	return 0
}
