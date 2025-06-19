package internal

import (
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"time"
)

func afterDur(then int64) string {
	if then == 0 {
		return "-"
	}
	now := time.Now().Unix()
	diff := then - now
	if diff == 0 {
		return "now"
	}

	units := []string{"year", "month", "day", "hour", "minute", "second"}
	unitSecs := []int64{31536000, 2592000, 86400, 3600, 60, 1}
	for i, unit := range units {
		if diff >= unitSecs[i] || unit == "seconds" {
			mults := math.Round(float64(diff) / float64(unitSecs[i]))
			if mults > 1 {
				unit += "s"
			}
			return fmt.Sprintf("in %v %v", mults, unit)
		}
	}
	return "Something went wrong!"
}

func processReq(req *http.Request, appData AppData) (bool, []byte) {
	req.Header.Set("X-API-Key", appData.Config.APIKey)

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln("Error sending request!")
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("Error reading response!")
	}
	statusOk := resp.StatusCode >= 200 && resp.StatusCode < 300
	return statusOk, body
}

func expiryString(t int64) string {
	expiry := "never"
	if t > 0 {
		expiry = time.Unix(t, 0).String()
	}
	return expiry
}
