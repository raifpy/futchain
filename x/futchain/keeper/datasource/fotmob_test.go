package datasource

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestFetch(t *testing.T) {

	var headers = make(http.Header)
	for key, value := range map[string]string{
		"accept": "*/*",
		"accept-language":    "en-US,en;q=0.7",
		"cache-control":      "no-cache",
		"dnt":                "1",
		"pragma":             "no-cache",
		"priority":           "u=1, i",
		"referer":            "https://www.fotmob.com/",
		"sec-ch-ua":          `"Chromium";v="140", "Not=A?Brand";v="24", "Brave";v="140"`,
		"sec-ch-ua-mobile":   "?0",
		"sec-ch-ua-platform": `"macOS"`,
		"sec-fetch-dest":     "empty",
		"sec-fetch-mode":     "cors",
		"sec-fetch-site":     "same-origin",
		"sec-gpc":            "1",
		"user-agent":         "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/140.0.0.0 Safari/537.36",
		"x-mas":              "eyJib2R5Ijp7InVybCI6Ii9hcGkvZGF0YS9tYXRjaGVzP2RhdGU9MjAyNTA5MDYmdGltZXpvbmU9RXVyb3BlJTJGSXN0YW5idWwmY2NvZGUzPUdCUiIsImNvZGUiOjE3NTcxNjIzMjEzNTEsImZvbyI6InByb2R1Y3Rpb246ZTUyYTNmYzE5Y2Y0YmY0NTY3ZTBlMzA3N2Q1OWQzNjVhNGEyYjNkNiJ9LCJzaWduYXR1cmUiOiIwOUVDNUVBNTNDQkNDNEFCODU3M0ZFMDQxN0REMTUwRCJ9"} {
		headers.Set(key, value)
	}

	ds := DatasourceFM{
		Client:  &http.Client{},
		BaseURL: "https://www.fotmob.com",
		Headers: headers,
	}

	tm, _ := context.WithTimeout(context.TODO(), time.Second*5)
	league, err := ds.Fetch(tm, WithTimezone("Europe/Istanbul"))
	if err != nil {
		t.Fatalf("failed to fetch data: %v", err)
	}

	for _, l := range league {
		for _, m := range l.Matches {
			fmt.Println(m)
		}
	}
}
