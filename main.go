package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/influxdb/influxdb/client"
)

type IPX800 struct {
	Day     string  `xml:"day"`
	Time0   string  `xml:"time0"`
	Analog0 float64 `xml:"analog0"`
	Analog1 float64 `xml:"analog1"`
	Analog2 float64 `xml:"analog2"`
	Analog3 float64 `xml:"analog3"`
}

const (
	MyIPXHost           = "http://files.pierrezemb.fr/status.xml"
	InfluxDBHost        = "localhost"
	InfluxDBPort        = 8086
	InfluxDB            = "celadon"
	InfluxDBMeasurement = "shapes"
	Time_Layout         = "02/01/2006 15:04:05"
)

func main() {

	log.Println("Starting IPX800 Watchdogs")

	log.Println("Connecting to InfluxDB")

	u, err := url.Parse(fmt.Sprintf("%s:%d", InfluxDBHost, InfluxDBPort))
	if err != nil {
		log.Fatal(err)
	}

	conf := client.Config{
		URL:      *u,
		Username: os.Getenv("INFLUX_USER"),
		Password: os.Getenv("INFLUX_PWD"),
	}

	con, err := client.NewClient(conf)
	if err != nil {
		log.Fatal(err)
	}

	dur, ver, err := con.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Happy as a Hippo! %v, %s", dur, ver)

	// based on https://stackoverflow.com/questions/16466320/is-there-a-way-to-do-repetitive-tasks-at-intervals-in-golang
	ticker := time.NewTicker(1 * time.Second)
	quit := make(chan struct{})

	for {
		select {
		case <-ticker.C:
			// do stuff
			go get_Data(con)
		case <-quit:
			ticker.Stop()
			return
		}
	}

	log.Println("Bye bye!")
}

func get_Data(con *client.Client) {

	var (
		ipx IPX800
		pts = make([]client.Point, 5)
	)

	// Get XML from IPX
	response, err := http.Get(MyIPXHost)
	if err != nil {
		log.Println(err)
		return
	}

	// Parsing XML
	if err = xml.NewDecoder(response.Body).Decode(&ipx); err != nil {
		log.Println(err)
		return
	}

	// Formatting TS
	// https://gobyexample.com/time-formatting-parsing
	// Mon Jan _2 15:04:05 2006
	t, _ := time.Parse(Time_Layout, ipx.Day+" "+ipx.Time0)
	if err != nil {
		log.Println(err)
		return
	}

	// Creating points
	pts[0] = client.Point{
		Name: "ipx800.temp",
		Tags: map[string]string{
			"ip":       MyIPXHost,
			"location": "cabine",
			"unit":     "°C",
		},
		Fields: map[string]interface{}{
			"value": strconv.FormatFloat((ipx.Analog0*0.323 - 50), 'f', 6, 64),
		},
		Timestamp: t,
		Precision: "s",
	}

	pts[1] = client.Point{
		Name: "ipx800.temp",
		Tags: map[string]string{
			"ip":       MyIPXHost,
			"location": "cabine",
			"unit":     "°C",
		},
		Fields: map[string]interface{}{
			"value": strconv.FormatFloat((ipx.Analog1*0.323 - 50), 'f', 6, 64),
		},
		Timestamp: t,
		Precision: "s",
	}

	pts[2] = client.Point{
		Name: "ipx800.intensity",
		Tags: map[string]string{
			"ip":   MyIPXHost,
			"unit": "A",
		},
		Fields: map[string]interface{}{
			"value": strconv.FormatFloat((ipx.Analog2 * 0.00323 * 5.52462), 'f', 6, 64),
		},
		Timestamp: t,
		Precision: "s",
	}

	pts[3] = client.Point{
		Name: "ipx800.voltage",
		Tags: map[string]string{
			"ip":   MyIPXHost,
			"unit": "A",
		},
		Fields: map[string]interface{}{
			"value": strconv.FormatFloat((ipx.Analog3  - 503) / 20), 'f', 6, 64),
		},
		Timestamp: t,
		Precision: "s",
	}
	pts[4] = client.Point{
		Name: "ipx800.power",
		Tags: map[string]string{
			"ip":   MyIPXHost,
			"unit": "V",
		},
		Fields: map[string]interface{}{
			"value": strconv.FormatFloat(((ipx.Analog2 * 0.00323 * 5.52462) * ((ipx.Analog3 - 503) / 20)), 'f', 6, 64),
		},
		Timestamp: t,
		Precision: "s",
	}

	// Pushing points into server
	bps := client.BatchPoints{
		Points:          pts,
		Database:        InfluxDB,
		RetentionPolicy: "default",
	}
	_, err = con.Write(bps)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(t)

}
