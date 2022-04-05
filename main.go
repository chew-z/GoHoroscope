package main

import (
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
	// loc = os.Getenv("LOCATION")
	location, _ = time.LoadLocation(city)
	// Point to where Swiss Ephem files are located on your system
	// It is a good practice to do it as initialization
	// even when not using files
	swephgo.SetEphePath([]byte("/usr/local/share/sweph/ephe"))
}

func main() {
	// Where the magic happens
	http.HandleFunc("/", CloudCharts)
	http.ListenAndServe(":8089", nil)

	// start := time.Now().UTC() // Start now
	// tx := julian(start)
	// fmt.Printf("%f\n", *tx)
	// rD := christian(tx)
	// rT := chrisToLocal(&rD)
	// fmt.Printf("%s\n", rT)

	swephgo.Close()
}
