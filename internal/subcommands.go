// SPDX-FileCopyrightText: 2025 Sayantan Santra <sayantan.santra689@gmail.com>
// SPDX-License-Identifier: MIT

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

type BackendConfig struct {
	Version        string `json:"version"`
	SiteURL        string `json:"site_url"`
	CapitalLetters bool   `json:"allow_capital_letters"`
	PublicMode     bool   `json:"public_mode"`
	PMExpiryDelay  int64  `json:"public_mode_expiry_delay"`
	SlugStyle      string `json:"slug_style"`
	SlugLength     int64  `json:"slug_length"`
	TryLongerSlug  bool   `json:"try_longer_slug"`
}

type JSONError struct {
	Reason string `json:"reason"`
}

func CreateLink(appData AppData) {
	log.SetFlags(0)
	payLoad := fmt.Sprintf(`{"shortlink":"%v","longlink":"%v","expiry_delay":%v}`, appData.Input2, appData.Input1, appData.Input3)
	req, _ := http.NewRequest("POST", appData.Config.URL+"/api/new", bytes.NewBufferString(payLoad))
	req.Header.Set("Content-Type", "application/json")

	ok, body := processReq(req, appData.Config)
	// This is for password based login
	// We kinda need to do some extra steps
	if appData.Config.APIKey == "" {
		if ok {
			fmt.Println("Shortlink:", string(body))
			fmt.Println("Expiry time is not reported for password based link creation.")
			fmt.Println("Run getall for that information.")
		} else {
			log.Fatalln(string(body))
		}
	} else {
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
}

func DeleteLink(appData AppData) {
	log.SetFlags(0)
	if appData.Input2 != "" {
		log.Fatalln("Too many arguments! Please see help.")
	}

	req, _ := http.NewRequest("DELETE", appData.Config.URL+"/api/del/"+appData.Input1, nil)
	ok, body := processReq(req, appData.Config)
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
	if appData.Config.APIKey == "" {
		log.Fatalln("The expand subcommand only works with an API key.")
	}

	req, _ := http.NewRequest("POST", appData.Config.URL+"/api/expand", bytes.NewBufferString(appData.Input1))
	ok, body := processReq(req, appData.Config)
	if ok {
		writer := tabwriter.NewWriter(os.Stdout, 0, 4, 4, ' ', 0)
		defer writer.Flush()
		var entry ExpandedURL
		json.Unmarshal(body, &entry)
		printpos := 0
		for printpos < len(entry.LongURL) {
			printlen := min(len(entry.LongURL)-printpos, 80)
			if printpos == 0 {
				fmt.Fprintln(writer, "Longlink:\t", entry.LongURL[printpos:printlen])
			} else {
				fmt.Fprintln(writer, "\t", entry.LongURL[printpos:printpos+printlen])
			}
			printpos += printlen
		}
		fmt.Fprintln(writer, "Hits:\t", entry.Hits)
		fmt.Fprintln(writer, "Expiry:\t", expiryString(entry.ExpiryTime))
	} else {
		var err JSONError
		json.Unmarshal(body, &err)
		log.Fatalln(err.Reason)
	}
}

func GetConfig(appData AppData) {
	log.SetFlags(0)
	if appData.Input1 != "" {
		log.Fatalln("Too many arguments! Please see help.")
	}

	req, _ := http.NewRequest("GET", appData.Config.URL+"/api/getconfig", nil)
	ok, body := processReq(req, appData.Config)
	if ok {
		var entry BackendConfig
		json.Unmarshal(body, &entry)
		fmt.Println("Version: ", entry.Version)
		if entry.SiteURL != "" {
			fmt.Println("Site URL: ", entry.SiteURL)
		}
		fmt.Println("Allow Capital Letters: ", entry.CapitalLetters)
		fmt.Println("Public Mode: ", entry.PublicMode)
		if entry.PMExpiryDelay > 0 {
			fmt.Println("Public Mode Expiry Delay: ", entry.PMExpiryDelay)
		}
		fmt.Println("Slug Style: ", entry.SlugStyle)
		if entry.SlugStyle == "UID" {
			fmt.Println("Slug Length: ", entry.SlugLength)
			fmt.Println("Try Longer Slug: ", entry.TryLongerSlug)
		}
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
	ok, body := processReq(req, appData.Config)
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
		defer writer.Flush()
		fmt.Fprintf(writer, "Short URL\tLong URL\tHits\tExpiry\n")
		fmt.Fprintf(writer, "---------\t--------\t----\t------\n")
		for _, entry := range slices.Backward(entries) {
			printpos := 0
			for printpos < len(entry.LongURL) {
				printlen := min(len(entry.LongURL)-printpos, 80)
				if printpos == 0 {
					fmt.Fprintf(writer, "%v\t%v\t%v\t%v\n", entry.ShortURL, entry.LongURL[printpos:printlen], entry.Hits, afterDur(entry.ExpiryTime))
				} else {
					fmt.Fprintf(writer, "\t%v\t\t\n", entry.LongURL[printpos:printpos+printlen])
				}
				printpos += printlen
			}
		}
	}
}
