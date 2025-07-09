// SPDX-FileCopyrightText: 2025 Sayantan Santra <sayantan.santra689@gmail.com>
// SPDX-License-Identifier: MIT

package internal

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"net/http/cookiejar"
	"os"
	"os/exec"
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

func doPasswordLogin(client *http.Client, config Config) {
	if config.Password == "" {
		fmt.Print("Empty password was supplied! Trying anyway.")
	}
	req, _ := http.NewRequest("POST", config.URL+"/api/login", bytes.NewBufferString(config.Password))
	jar, _ := cookiejar.New(nil)
	client.Jar = jar

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln("Error trying to log in!")
	}
	defer resp.Body.Close()

	if err != nil {
		log.Fatalln("Error reading response from login route!")
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Fatalln("Error logging in! Check your password.")
	}
}

func processReq(req *http.Request, config Config) (bool, []byte) {
	client := http.DefaultClient
	if config.APIKey == "" {
		doPasswordLogin(client, config)
	} else {
		req.Header.Set("X-API-Key", config.APIKey)
	}

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

func readPass() string {
	// Disable printing to terminal and handle all characters
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
	var rs []rune
	inp := bufio.NewReader(os.Stdin)
	for {
		r, _, err := inp.ReadRune()
		if err != nil {
			panic(err)
		}
		if r == '\n' {
			fmt.Print("\n")
			break
		} else if r == '\u007f' { // backspace
			if len(rs) > 0 {
				fmt.Printf("\033[1D\033[K")
				rs = rs[:len(rs)-1]
			}
			continue
		} else {
			fmt.Print("*")
			rs = append(rs, r)
		}
	}
	return string(rs)
}
