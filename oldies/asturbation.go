package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
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
	// birthDay, birthMonth, birthYear int
	// birthHour, latitude, longitude  float64
	location string
)

func init() {
	location = os.Getenv("LOCATION")
	// Point to where Swiss Ephem files are located on your system
	swephgo.SetEphePath([]byte("/usr/local/share/sweph/ephe"))
}

func main() {
	http.HandleFunc("/", httpserver)
	http.ListenAndServe(":8089", nil)

	swephgo.Close()
}

func httpserver(w http.ResponseWriter, _ *http.Request) {
	var line3DColor = []string{
		"#313695", "#4575b4", "#74add1", "#abd9e9", "#e0f3f8",
		"#fee090", "#fdae61", "#f46d43", "#d73027", "#a50026",
	}
	page := components.NewPage()

	// moon := charts.NewLine3D()
	mercury := charts.NewLine3D()
	// venus := charts.NewLine3D()
	// mars := charts.NewLine3D()

	// moon.SetGlobalOptions(
	// 	charts.WithTitleOpts(opts.Title{Title: "Moon"}),
	// 	charts.WithVisualMapOpts(opts.VisualMap{
	// 		Calculable: true,
	// 		Max:        1,
	// 		InRange:    &opts.VisualMapInRange{Color: line3DColor},
	// 	}),
	// 	charts.WithGrid3DOpts(opts.Grid3D{
	// 		ViewControl: &opts.ViewControl{
	// 			AutoRotate: true,
	// 		},
	// 	}),
	// )
	mercury.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "Mercury"}),
		charts.WithVisualMapOpts(opts.VisualMap{
			Calculable: true,
			Max:        1,
			InRange:    &opts.VisualMapInRange{Color: line3DColor},
		}),
	)
	// 	charts.WithGrid3DOpts(opts.Grid3D{
	// 		ViewControl: &opts.ViewControl{
	// 			AutoRotate: true,
	// 		},
	// 	}),
	// )
	// venus.SetGlobalOptions(
	// 	charts.WithTitleOpts(opts.Title{Title: "Venus"}),
	// 	charts.WithVisualMapOpts(opts.VisualMap{
	// 		Calculable: true,
	// 		Max:        1,
	// 		InRange:    &opts.VisualMapInRange{Color: line3DColor},
	// 	}),
	// 	charts.WithGrid3DOpts(opts.Grid3D{
	// 		ViewControl: &opts.ViewControl{
	// 			AutoRotate: true,
	// 		},
	// 	}),
	// )
	// mars.SetGlobalOptions(
	// 	charts.WithTitleOpts(opts.Title{Title: "Mars"}),
	// 	charts.WithVisualMapOpts(opts.VisualMap{
	// 		Calculable: true,
	// 		Max:        1,
	// 		InRange:    &opts.VisualMapInRange{Color: line3DColor},
	// 	}),
	// 	charts.WithGrid3DOpts(opts.Grid3D{
	// 		ViewControl: &opts.ViewControl{
	// 			AutoRotate: true,
	// 		},
	// 	}),
	// )

	// moon.AddSeries("moon", generateData(swephgo.SeMoon, time.Now().UTC(), 10, 6))
	mercury.AddSeries("mercury", generateData(swephgo.SeMercury, time.Now().UTC(), 10, 6))
	// venus.AddSeries("venus", generateData(swephgo.SeVenus, time.Now().UTC(), 1, 6))
	// mars.AddSeries("mars", generateData(swephgo.SeMars, time.Now().UTC(), 1, 6))
	page.AddCharts(
		mercury,
		// moon,
		// venus,
		// mars,
	)
	page.Render(w)
}

func generateData(ipl int, startTime time.Time, years int, months int) []opts.Chart3DData {
	data := make([][3]float64, 0)
	start := nod(startTime)
	end := start.AddDate(years, months, 0)
	for d := start; d.After(end) == false; d = d.AddDate(0, 0, 1) {
		julianDay := swephgo.Julday(d.Year(), int(d.Month()), d.Day(), float64(d.Hour()), swephgo.SeGregCal)
		p, _ := phase(&julianDay, ipl)
		waldo, _ := waldo(&julianDay, ipl)
		data = append(data, [3]float64{waldo[2], waldo[5], p})
	}
	ret := make([]opts.Chart3DData, 0, len(data))
	for _, d := range data {
		ret = append(ret, opts.Chart3DData{Value: []interface{}{d[0], d[1], d[2]}})
	}
	return ret
}

/*
What is a phase (ilumination) of a planet?
https://groups.io/g/swisseph/message/7327
*/
func phase(julianDay *float64, planet int) (float64, error) {
	iflag := swephgo.SeflgSwieph // use SWISSEPH ephemeris, default
	attr := make([]float64, 20)
	serr := make([]byte, 256)
	eclflag := swephgo.PhenoUt(*julianDay, planet, iflag, attr, serr)
	if eclflag == swephgo.Err {
		log.Printf("Error %d %s", eclflag, string(serr))
		return 0.0, errors.New(string(serr))
	}
	return attr[1], nil
}

/*
Where is a planet (longitude, latitude, distance, speed in long., speed in lat., and speed in dist.)
*/
func waldo(julianDay *float64, planet int) ([]float64, error) {
	iflag := swephgo.SeflgSwieph + swephgo.SeflgSpeed // use SWISSEPH ephemeris, default and calculate speed of bodies
	x2 := make([]float64, 6)
	serr := make([]byte, 256)
	eclflag := swephgo.CalcUt(*julianDay, planet, iflag, x2, serr)
	if eclflag == swephgo.Err {
		return x2, errors.New(string(serr))
	}
	return x2, nil
}

/* Julian date back to calendar date */
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

/* Begining of the Day */
func bod(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.Local)
}

/* Noon of the Day */
func nod(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 12, 0, 0, 0, time.Local)
}

func bool2int(b bool) int {
	if b {
		return 1
	}
	return 0
}
