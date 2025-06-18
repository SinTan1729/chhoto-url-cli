package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

func createLink(appData AppData) {
	log.SetFlags(0)
	payLoad := fmt.Sprintf(`{"shorturl":"%v","longurl":"%v","expiry_delay":%v}`, appData.Input1, appData.Input2, appData.Input3)
	req, _ := http.NewRequest("POST", appData.Config.URL+"/api/new", bytes.NewBufferString(payLoad))
	req.Header.Set("Content-Type", "application/json")
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
	fmt.Println(string(body))
}

func deleteLink(appData AppData) {
	log.SetFlags(0)
	log.Fatalln("Delete")
}

func expandLink(appData AppData) {
	log.SetFlags(0)
	log.Fatalln("Expand")
}

func getAll(appData AppData) {
	log.SetFlags(0)
	req, _ := http.NewRequest("GET", appData.Config.URL+"/api/all", nil)
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

	var entries []URLEntry
	err = json.Unmarshal(body, &entries)
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
