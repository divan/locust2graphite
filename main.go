package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	graph "github.com/marpaia/graphite-golang"
)

const (
	STATS_URI = "/stats/requests"
)

var (
	locust   = flag.String("locust", "http://10.51.110.41:8089", "Locust host")
	addr     = flag.String("host", "graphite.divan", "Graphite host")
	port     = flag.Int("port", 2003, "Graphite port")
	interval = flag.Duration("interval", 1*time.Second, "Interval for sending stats")

	hostname, _ = os.Hostname()
	statsBase   = fmt.Sprintf("stats.%s.rest_api.", hostname)
	graphite    *graph.Graphite
)

func sendToGraphite(stats LocustStats) {
	fmt.Print(".")
	graphite.SimpleSend(statsBase+"fail_ratio", stats.FailRatio)
	totalStat, err := extractTotalStat(stats)
	if err != nil {
		log.Println(err)
		return
	}
	graphite.SimpleSend(statsBase+"total.median_response_time", totalStat.MedianResponseTime)
	graphite.SimpleSend(statsBase+"total.avg_response_time", totalStat.AvgResponseTime)
	graphite.SimpleSend(statsBase+"total.current_rps", totalStat.CurrentRps)
	graphite.SimpleSend(statsBase+"total.num_failures", totalStat.NumFailures)
}

func extractTotalStat(stats LocustStats) (LocustStat, error) {
	for _, stat := range stats.Stats {
		if stat.Name == "Total" {
			return stat, nil
		}
	}
	return LocustStat{}, fmt.Errorf("Cannot find Total stats in Locust requests stats...")
}

func getStats() (LocustStats, error) {
	var ret LocustStats

	resp, err := http.Get(*locust + STATS_URI)
	if err != nil {
		return ret, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return ret, err
	}

	err = json.Unmarshal(body, &ret)
	if err != nil {
		return ret, err
	}

	return ret, nil
}

func main() {
	flag.Parse()
	log.Printf("Using graphite on %s:%d\n", *addr, *port)
	var err error
	graphite, err = graph.NewGraphite(*addr, *port)
	if err != nil {
		log.Fatal(err)
	}
	for {
		stats, err := getStats()
		if err != nil {
			log.Println("Error reading stats from Locust:", err)
			time.Sleep(*interval)
			continue
		}
		sendToGraphite(stats)
		time.Sleep(*interval)
	}
}
