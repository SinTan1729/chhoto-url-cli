// SPDX-FileCopyrightText: 2025 Sayantan Santra <sayantan.santra689@gmail.com>
// SPDX-License-Identifier: MIT

package main

import (
	"github.com/SinTan1729/chhoto-url-cli/internal"
	"log"
)

func main() {
	log.SetFlags(0)
	appData := internal.ParseData()

	switch appData.Subcommand {
	case "new":
		internal.CreateLink(appData)
	case "delete":
		internal.DeleteLink(appData)
	case "expand":
		internal.ExpandLink(appData)
	case "getall":
		internal.GetAll(appData)
	default:
		log.Fatalln(appData.Subcommand, "is not a valid subcommand. Please see help.")
	}
}
