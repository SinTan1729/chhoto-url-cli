// SPDX-FileCopyrightText: 2025 Sayantan Santra <sayantan.santra689@gmail.com>
// SPDX-License-Identifier: MIT

package main

import (
	"log"
)

func main() {
	log.SetFlags(0)
	appData := parseData()

	switch appData.Subcommand {
	case "new":
		createLink(appData)
	case "delete":
		deleteLink(appData)
	case "expand":
		expandLink(appData)
	case "getall":
		getAll(appData)
	default:
		log.Fatalln(appData.Subcommand, "is not a valid subcommand. Please see help.")
	}
}
