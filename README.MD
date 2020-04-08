## What is this
This is console app for measuring pages loading speed. 
It is based on Google PageSpeedOnline API https://developers.google.com/speed/docs/insights/v5/get-started.
 
You can pass file with list of sites to measure and receive 
some parameters for each of it in YAML format.

#### What parameters you can receive from it

* `firstMeaningfulPaint` (sec) measures when primary content of page is visible. More: https://web.dev/first-meaningful-paint
* `timeToInteractive` (sec) measures when page is becoming fully interactive. More: https://web.dev/interactive
* `speedIndex` (0-100 points) measures how quickly content is displayed. More: https://web.dev/speed-index/
* `firstContentfulPaint` (sec) measures when first piece of DOM is rendered. More: https://web.dev/first-contentful-paint/

## Inapp help:
```
Usage of SiteSpeed:
     -apiKey string
           (required)
           API key for using Google PageSpeedOnline. It can be received at 
           https://developers.google.com/speed/docs/insights/v5/get-started#key.
     -duration string
           Duration to sleep between requests for each site. Default value is 1.5s
           Minimum time between requests to pagespeedonline api is 1s, but even default 1.5s is too little 
           sometimes. A duration string is a possibly signed sequence of decimal numbers, each with optional 
           fraction and a unit suffix, such as "300ms", "1.5h" or "2h45m".
           Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h". (default "1.5s")
     -noLogs
           If set, there is no log output to STDOUT.
     -output string
           Path to output in YAML format. Default value is "".
           If it is empty (default value), then output to STDOUT.
     -sites string
           (required) 
           Path to file with sites. On each row there must be only one site, row should end with comma (,) and
           can be surrounded with double quotes (").
```

## How to use it
You can use command `go build` and then run resulting binary.