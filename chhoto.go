// SPDX-FileCopyrightText: 2025 Sayantan Santra <sayantan.santra689@gmail.com>
// SPDX-License-Identifier: MIT

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func main() {
	// GET a JSON object
	id := 12
	var post placeholder
	resp, err := http.Get(fmt.Sprintf("https://jsonplaceholder.typicode.com/posts/%d", id))
	if err != nil {
		fmt.Println("could not connect to jsonplaceholder.typicode.com:", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	json.Unmarshal(body, &post)
	fmt.Println(post.Title)

	config := LoadConfig()
	fmt.Println(config.URL)
	fmt.Println(config.APIKey)
}

type placeholder struct {
	ID     int    `json:"id,omitempty"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	UserID int    `json:"userId"`
}
