// SPDX-FileCopyrightText: 2025 Sayantan Santra <sayantan.santra689@gmail.com>
// SPDX-License-Identifier: MIT

package main

import (
	"log"
)

func main() {
	log.SetFlags(0)
	appData := ParseData()

	switch appData.Subcommand {
	case "new":
		CreateLink(appData)
	case "delete":
		DeleteLink(appData)
	case "expand":
		ExpandLink(appData)
	case "getall":
		GetAll(appData)
	default:
		log.Fatalln(appData.Subcommand, "is not a valid subcommand. Please see help.")
	}
}

type placeholder struct {
	ID     int    `json:"id,omitempty"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	UserID int    `json:"userId"`
}
