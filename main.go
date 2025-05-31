package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/fatih/color"
)

type Hourly struct {
	Time          []string  `json:"time"`
	Temperature2m []float64 `json:"temperature_2m"`
}

type Response struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timezone  string  `json:"timezone"`
	Hourly    Hourly  `json:"hourly"`
}

func formatTemp(temp float64) string {
	var c *color.Color

	switch {
	case temp <= -10:
		c = color.New(color.FgBlue)
	case temp <= 5:
		c = color.New(color.FgHiBlue)
	case temp <= 10:
		c = color.New(color.FgHiCyan)
	case temp <= 15:
		c = color.New(color.FgHiGreen)
	case temp <= 20:
		c = color.New(color.FgHiYellow)
	case temp <= 25:
		c = color.New(color.FgHiRed)
	case temp <= 30:
		c = color.New(color.FgRed)
	default:
		c = color.New(color.FgWhite)
	}

	// space for alignment
	formatted := fmt.Sprintf("%.1f", temp)
	if temp < 10 && temp > -10 {
		formatted = " " + formatted
	}

	return c.Sprintf("%s°C", formatted)
}

func getCoords() (float64, float64, string, string, error) {
	resp, err := http.Get("https://ipwho.is/")

	if err != nil {
		return 0, 0, "", "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, 0, "", "", err
	}

	var data struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		City      string  `json:"city"`
		Country   string  `json:"country"`
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return 0, 0, "", "", err
	}

	return data.Latitude, data.Longitude, data.City, data.Country, nil
}

func main() {
	lat, lon, city, country, err := getCoords()
	if err != nil {
		log.Fatalf("Get coords error: %v", err)
	}

	url := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&hourly=temperature_2m&timezone=auto", lat, lon)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Request error: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Read response error: %v", err)
	}

	var data Response
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Fatalf("JSON parse error: %v", err)
	}

	layout := "2006-01-02T15:04"
	daily := make(map[string][]float64)

	for i, iso := range data.Hourly.Time {
		t, err := time.Parse(layout, iso)
		if err != nil {
			log.Fatalf("Parse time error: %s: %v", iso, err)
		}
		date := t.Format("(Mon) Jan 02, 2006")
		temp := data.Hourly.Temperature2m[i]
		daily[date] = append(daily[date], temp)
	}

	// Current temperature
	if len(data.Hourly.Temperature2m) > 0 {
		now := time.Now()
		hour := now.Hour()
		currentTemp := data.Hourly.Temperature2m[hour]
		fmt.Printf("\n%s %s\n\n", country, "/ "+city+": "+formatTemp(currentTemp))
	}

	layoutDate := "(Mon) Jan 02, 2006"

	// Parse dates
	parsedDates := make([]time.Time, 0, len(daily))
	dateMap := make(map[time.Time]string)

	for dateStr := range daily {
		t, err := time.Parse(layoutDate, dateStr)
		if err != nil {
			log.Fatalf("Parse date error: %s: %v", dateStr, err)
		}
		parsedDates = append(parsedDates, t)
		dateMap[t] = dateStr
	}

	// Sort dates
	sort.Slice(parsedDates, func(i, j int) bool {
		return parsedDates[i].Before(parsedDates[j])
	})

	// Print dates
	for _, t := range parsedDates {
		dateStr := dateMap[t]
		temps := daily[dateStr]
		min, max := temps[0], temps[0]
		for _, temp := range temps[1:] {
			if temp < min {
				min = temp
			}
			if temp > max {
				max = temp
			}
		}
		fmt.Printf("%s - %s : %s\n", formatTemp(min), formatTemp(max), dateStr)
	}
}
