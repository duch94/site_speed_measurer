package main

import "fmt"

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
