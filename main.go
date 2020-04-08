package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type SiteSpeed struct {
	date                 string
	url                  string
	firstMeaningfulPaint float64
	timeToInteractive    float64
	speedIndex           float64
	firstContentfulPaint float64
}

func (ss *SiteSpeed) GetYaml() string {
	res := fmt.Sprintf(`- date: %s
  url: %s
  firstMeaningfulPaint: %f
  timeToInteractive: %f
  speedIndex: %f
  firstContentfulPaint: %f`,
		ss.date, ss.url, ss.firstMeaningfulPaint, ss.timeToInteractive, ss.speedIndex, ss.firstContentfulPaint)
	return res
}

func main() {
	filePathHelp := `(required) 
Path to file with sites. On each row there must be only one site, row should end with comma (,) and
can be surrounded with double quotes (").`
	filePathPtr := flag.String("sites", "", filePathHelp)
	apiKeyHelp := `(required)
API key for using Google PageSpeedOnline. It can be received at 
https://developers.google.com/speed/docs/insights/v5/get-started#key.`
	apiKeyPtr := flag.String("apiKey", "", apiKeyHelp)
	durationHelp := `Duration to sleep between requests for each site. Default value is 1.5s
Minimum time between requests to pagespeedonline api is 1s, but even default 1.5s is too little 
sometimes. A duration string is a possibly signed sequence of decimal numbers, each with optional 
fraction and a unit suffix, such as "300ms", "1.5h" or "2h45m".
Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".`
	unparsedDurPtr := flag.String("duration", "1.5s", durationHelp)
	outputHelp := `Path to output in YAML format. Default value is "".
If it is empty (default value), then output to STDOUT.`
	outputPathPtr := flag.String("output", "", outputHelp)
	noLogsHelp := `If set, there is no log output to STDOUT.`
	noLogsPtr := flag.Bool("noLogs", false, noLogsHelp)
	flag.Parse()

	if *noLogsPtr {
		log.SetOutput(ioutil.Discard)
	}

	dur, err := time.ParseDuration(*unparsedDurPtr)
	if err != nil {
		panic(err)
	}

	sites := loadSites(*filePathPtr)
	ch := make(chan SiteSpeed)
	for i := range sites {
		go measureLoadingSpeed(sites[i], ch, apiKeyPtr)
		time.Sleep(dur)
	}
	var speeds []SiteSpeed
	for range sites {
		speeds = append(speeds, <-ch)
	}

	outputSiteSpeeds(speeds, *outputPathPtr)
}

func outputSiteSpeeds(siteSpeeds []SiteSpeed, output string) {
	outputString := "---\n"
	for i := range siteSpeeds {
		outputString += siteSpeeds[i].GetYaml() + "\n\n"
	}
	if len(output) == 0 {
		fmt.Println(outputString)
	} else {
		err := ioutil.WriteFile(output, []byte(outputString), 0644)
		if err != nil {
			panic(err)
		}
	}
}

func loadSites(filePath string) []string {
	sitesBytes, err := ioutil.ReadFile(filePath)
	check(err)
	var sites []string
	sites = strings.Split(string(sitesBytes), ",\n")
	for i := range sites {
		sites[i] = strings.ReplaceAll(sites[i], "\"", "")
		sites[i] = strings.ReplaceAll(sites[i], "\n", "")
		sites[i] = strings.ReplaceAll(sites[i], "\r", "")
	}
	return sites
}

func replaceTrash(s string) string {
	return strings.ReplaceAll(s, "\u00a0s", "")
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func parseSiteSpeed(jsonToParse map[string]map[string]map[string]map[string]string, url string) SiteSpeed {
	fmpStr := jsonToParse["lighthouseResult"]["audits"]["first-meaningful-paint"]["displayValue"]
	fmpStr = replaceTrash(fmpStr)
	fmp, err := strconv.ParseFloat(fmpStr, 64)
	if err != nil {
		log.Printf("Error while parsing firstMeaningfulPaing for url=%s: assinged 0", url)
	}

	ttiStr := jsonToParse["lighthouseResult"]["audits"]["interactive"]["displayValue"]
	ttiStr = replaceTrash(ttiStr)
	tti, err := strconv.ParseFloat(ttiStr, 64)
	if err != nil {
		log.Printf("Error while parsing timeToInteractive for url=%s: assinged 0", url)
	}

	siStr := jsonToParse["lighthouseResult"]["audits"]["speed-index"]["displayValue"]
	siStr = replaceTrash(siStr)
	si, err := strconv.ParseFloat(siStr, 64)
	if err != nil {
		log.Printf("Error while parsing speedIndex for url=%s: assinged 0", url)
	}

	fcpStr := jsonToParse["lighthouseResult"]["audits"]["first-contentful-paint"]["displayValue"]
	fcpStr = replaceTrash(fcpStr)
	fcp, err := strconv.ParseFloat(fcpStr, 64)
	if err != nil {
		log.Printf("Error while parsing firstContentfulPaint for url=%s: assinged 0", url)
	}

	y, m, d := time.Now().Date()
	today := fmt.Sprintf("%d-%d-%d", y, m, d)

	return SiteSpeed{
		date:                 today,
		url:                  url,
		firstMeaningfulPaint: fmp,
		timeToInteractive:    tti,
		speedIndex:           si,
		firstContentfulPaint: fcp,
	}
}

func measureLoadingSpeed(url string, ch chan SiteSpeed, apiKeyPtr *string) {
	pageSpeedUrl := fmt.Sprintf("https://www.googleapis.com/pagespeedonline/v5/runPagespeed?key=%s&url=%s",
		*apiKeyPtr, url)
	log.Printf("Measuring loading speed for url=%s", url)
	resp, err := http.Get(pageSpeedUrl)
	log.Printf("Performed GET request to url=%s", pageSpeedUrl)
	if err != nil {
		log.Printf("Error occured while getting results for url=%s, retrying...", url)
		time.Sleep(1)
		resp, err = http.Get(pageSpeedUrl)
		if err != nil {
			log.Panicf("Error occured while getting results: %s", err)
		}
	}

	defer resp.Body.Close()
	var t map[string]map[string]map[string]map[string]string
	_ = json.NewDecoder(resp.Body).Decode(&t)

	ss := parseSiteSpeed(t, url)
	ch <- ss
	log.Printf("Measured load speed for url=%s", url)
}
