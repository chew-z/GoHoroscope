package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/markcheno/go-talib"
	"github.com/mshafiee/swephgo"
)

/*Quotes ...
 */
type Quotes struct {
	DefaultChartInterval string          `json:"_default_chart_interval,omitempty"`
	RefPrice             float64         `json:"_ref_price,omitempty"`
	Bars                 [][]interface{} `json:"_d"`
}

/*Candle ...
 */
type Candle struct {
	Time string
	// BOSSA API is OHLC go-echarts OCLH
	OHLC [4]float64
}

var (
	apiURL      = os.Getenv("API_URL")
	asset       = os.Getenv("ASSET")
	city        = os.Getenv("CITY")
	location, _ = time.LoadLocation(city)
	userAgent   = randUserAgent()
	kd          [100]Candle
	indicator0  [100]float64
	indicator1  [100]float64
	indicator2  [100]float64
	indicator3  [100]float64
	indicator4  [100]float64
	typical     [100]float64
	moon        [100]float64
	mercury     [100]float64
	venus       [100]float64
)

/*CloudCharts ...
 */
func CloudCharts(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	var high, low [100]float64
	timeframe := "D1"
	if a := query.Get("a"); a != "" {
		asset = a
	}
	if t := query.Get("t"); t != "" {
		timeframe = t
	}
	if quotes := getQuotes(asset, timeframe); quotes != nil {
		for i, bar := range quotes.Bars {
			var tmp Candle
			tm := int64(bar[0].(float64))
			time := time.Unix(0, tm*int64(time.Millisecond))
			tmp.Time = time.In(location).Format("Jan _2 15:04")
			o, _ := bar[1].(float64)
			h, _ := bar[2].(float64)
			l, _ := bar[3].(float64)
			c, _ := bar[4].(float64)
			high[i] = h
			low[i] = l
			typical[i] = (h + l + c) / 3.0    // typical price - HLC/3
			tmp.OHLC = [4]float64{o, c, l, h} // OHLC to OCLH
			kd[i] = tmp
			jD := swephgo.Julday(time.Year(), int(time.Month()), time.Day(), float64(time.Hour()), swephgo.SeGregCal)
			moon[i], _ = Phase(&jD, swephgo.SeMoon)
			mercury[i], _ = Phase(&jD, swephgo.SeMercury)
			venus[i], _ = Phase(&jD, swephgo.SeVenus)

		}
		ma0 := talib.Ma(typical[:], 10, talib.SMA)
		ma1 := talib.Ma(typical[:], 20, talib.SMA)
		copy(indicator0[:], ma0)
		copy(indicator1[:], ma1)
		copy(indicator2[:], moon[:])
		copy(indicator3[:], mercury[:])
		copy(indicator4[:], venus[:])
		bars := ohlcChart()
		indicators := indicatorsChart()
		page := components.NewPage()
		page.AddCharts(
			bars,
			indicators,
		)
		page.Render(w)
	} else {
		http.Error(w, "Something went wrong, can't render chart", http.StatusInternalServerError)
	}
}

func indicatorsChart() *charts.Line {
	lineChart := charts.NewLine()
	x := make([]string, 100)
	z := make([]opts.LineData, 100)
	m := make([]opts.LineData, 100)
	v := make([]opts.LineData, 100)
	for i := 0; i < len(kd); i++ {
		x[i] = kd[i].Time
		z[i] = opts.LineData{Value: indicator2[i]}
		m[i] = opts.LineData{Value: indicator3[i]}
		v[i] = opts.LineData{Value: indicator4[i]}
	}

	lineChart.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "Moon - Mercury - Venus",
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Scale: true,
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Type:       "slider",
			Start:      21,
			End:        100,
			XAxisIndex: []int{0},
		}),
		charts.WithTooltipOpts(opts.Tooltip{
			Show: true,
		}),
	)
	lineChart.SetXAxis(x).AddSeries("Moon", z)
	lineChart.SetXAxis(x).AddSeries("Mercury", m)
	lineChart.SetXAxis(x).AddSeries("Venus", v)
	return lineChart
}
func ohlcChart() *charts.Kline {
	// create a new chart instance
	kline := charts.NewKLine()
	x := make([]string, 100)
	y := make([]opts.KlineData, 100)
	v := make([]opts.LineData, 100)
	z := make([]opts.LineData, 100)
	for i := 0; i < len(kd); i++ {
		x[i] = kd[i].Time
		y[i] = opts.KlineData{Value: kd[i].OHLC}
		v[i] = opts.LineData{Value: indicator0[i]}
		z[i] = opts.LineData{Value: indicator1[i]}
	}
	kline.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{PageTitle: fmt.Sprintf("%s - %.5g", asset, kd[99].OHLC[1])}),
		charts.WithTitleOpts(opts.Title{
			Title:    fmt.Sprintf("%s - %.5g", asset, kd[99].OHLC[1]),
			Subtitle: fmt.Sprintf("%.5g - %.5g", kd[99].OHLC[2], kd[99].OHLC[3]),
		}),
		charts.WithXAxisOpts(opts.XAxis{
			SplitNumber: 20,
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Scale: true,
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Type:       "slider",
			Start:      25,
			End:        100,
			XAxisIndex: []int{0},
		}),
		charts.WithTooltipOpts(opts.Tooltip{
			Show: true,
		}),
		charts.WithToolboxOpts(opts.Toolbox{
			Show:  true,
			Right: "20%",
			Feature: &opts.ToolBoxFeature{
				SaveAsImage: &opts.ToolBoxFeatureSaveAsImage{
					Show:  true,
					Type:  "png",
					Title: "save as image",
				},
			}},
		),
	)
	kline.SetXAxis(x).AddSeries(asset, y).
		SetSeriesOptions(
			charts.WithMarkPointNameTypeItemOpts(opts.MarkPointNameTypeItem{
				Name:     "Maximum",
				Type:     "max",
				ValueDim: "highest",
			}),
			charts.WithMarkPointNameTypeItemOpts(opts.MarkPointNameTypeItem{
				Name:     "Minimum",
				Type:     "min",
				ValueDim: "lowest",
			}),
			charts.WithMarkPointStyleOpts(opts.MarkPointStyle{
				Symbol: []string{"pin"},
				Label: &opts.Label{
					Show: true,
				},
			}),
			charts.WithItemStyleOpts(opts.ItemStyle{
				Color:        "#00da3c",
				Color0:       "#ec0000",
				BorderColor:  "#008F28",
				BorderColor0: "#8A0000",
			}),
		)
	lineChart := charts.NewLine()
	lineChart.SetGlobalOptions(
		charts.WithYAxisOpts(opts.YAxis{
			Scale: true,
		}),
	)
	lineChart.SetXAxis(x).AddSeries("MA0", v).AddSeries("MA1", z)
	kline.Overlap(lineChart)
	return kline
}

func getQuotes(asset string, timeframe string) *Quotes {
	var quotes Quotes
	apiURL := fmt.Sprintf("%s%s.", apiURL, asset)
	if timeframe != "" {
		apiURL += "/" + timeframe
	}
	request, _ := http.NewRequest("GET", apiURL, nil)
	request.Header.Set("User-Agent", userAgent)
	if response, err := client.Do(request); err == nil {
		if err := json.NewDecoder(response.Body).Decode(&quotes); err != nil {
			log.Fatal(err)
		}
	}
	return &quotes
}
