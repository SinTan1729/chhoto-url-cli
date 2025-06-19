package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"slices"
	"text/tabwriter"
)

type URLEntry struct {
	ShortURL   string `json:"shortlink"`
	LongURL    string `json:"longlink"`
	Hits       int64  `json:"hits"`
	ExpiryTime int64  `json:"expiry_time"`
}

type ExpandedURL struct {
	ShortURL   string `json:"shorturl,omitempty"`
	LongURL    string `json:"longurl,omitempty"`
	Hits       int64  `json:"hits,omitempty"`
	ExpiryTime int64  `json:"expiry_time"`
}

type JSONError struct {
	Reason string `json:"reason"`
}

func CreateLink(appData AppData) {
	log.SetFlags(0)
	payLoad := fmt.Sprintf(`{"shortlink":"%v","longlink":"%v","expiry_delay":%v}`, appData.Input2, appData.Input1, appData.Input3)
	req, _ := http.NewRequest("POST", appData.Config.URL+"/api/new", bytes.NewBufferString(payLoad))
	req.Header.Set("Content-Type", "application/json")

	ok, body := processReq(req, appData)
	if ok {
		var entry ExpandedURL
		json.Unmarshal(body, &entry)
		fmt.Println("Shortlink: ", entry.ShortURL)
		fmt.Println("Expiry: ", expiryString(entry.ExpiryTime))
	} else {
		var err JSONError
		json.Unmarshal(body, &err)
		log.Fatalln(err.Reason)
	}
}

func DeleteLink(appData AppData) {
	log.SetFlags(0)
	if appData.Input2 != "" {
		log.Fatalln("Too many arguments! Please see help.")
	}

	req, _ := http.NewRequest("DELETE", appData.Config.URL+"/api/del/"+appData.Input1, nil)
	ok, body := processReq(req, appData)
	if ok {
		fmt.Printf("Shortlink %v was successfully deleted!", appData.Input1)
	} else {
		var err JSONError
		json.Unmarshal(body, &err)
		log.Fatalln(err.Reason)
	}
}

func ExpandLink(appData AppData) {
	log.SetFlags(0)
	if appData.Input2 != "" {
		log.Fatalln("Too many arguments! Please see help.")
	}

	req, _ := http.NewRequest("POST", appData.Config.URL+"/api/expand", bytes.NewBufferString(appData.Input1))
	ok, body := processReq(req, appData)
	if ok {
		var entry ExpandedURL
		json.Unmarshal(body, &entry)
		fmt.Println("Longlink: ", entry.LongURL)
		fmt.Println("Hits: ", entry.Hits)
		fmt.Println("Expiry: ", expiryString(entry.ExpiryTime))
	} else {
		var err JSONError
		json.Unmarshal(body, &err)
		log.Fatalln(err.Reason)
	}
}

func GetAll(appData AppData) {
	log.SetFlags(0)
	if appData.Input1 != "" {
		log.Fatalln("Too many arguments! Please see help.")
	}

	req, _ := http.NewRequest("GET", appData.Config.URL+"/api/all", nil)
	ok, body := processReq(req, appData)
	if !ok {
		log.Fatalln("Received error from the server!")
	}

	var entries []URLEntry
	err := json.Unmarshal(body, &entries)
	if err != nil {
		log.Fatalln("Error reading JSON!")
	}
	if len(entries) == 0 {
		fmt.Println("No links were returned.")
	} else {
		writer := tabwriter.NewWriter(os.Stdout, 0, 4, 4, ' ', 0)
		fmt.Fprintf(writer, "Short URL\tLong URL\tHits\tExpiry\n")
		fmt.Fprintf(writer, "---------\t--------\t----\t------\n")
		for _, entry := range slices.Backward(entries) {
			fmt.Fprintf(writer, "%v\t%v\t%v\t%v\n", entry.ShortURL, entry.LongURL, entry.Hits, afterDur(entry.ExpiryTime))
		}
		writer.Flush()
	}
}
