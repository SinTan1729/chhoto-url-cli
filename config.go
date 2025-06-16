package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"text/tabwriter"
)

type Config struct {
	URL    string `json:"url"`
	APIKey string `json:"apiKey"`
}

func LoadConfig() Config {
	configDir, ok := os.LookupEnv("XDG_CONFIG_HOME")
	if !ok {
		configDir = "~/.config"
	}
	file, err := os.Open(configDir + "/chhoto/config.json")
	if err != nil {
		fmt.Println("Could not load config from " + configDir + "/chhoto/config.json. Quitting!")
		os.Exit(1)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	config := Config{}
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Println("error:", err)
	}

	urlFlag := flag.String("url", "", "URL of the Chhoto URL server.")
	apiFlag := flag.String("api-key", "", "API Key of the Chhoto URL server.")
	flag.Usage = func() {
		writer := tabwriter.NewWriter(flag.CommandLine.Output(), 0, 4, 4, ' ', 0)
		fmt.Fprintf(writer, "Chhoto URL CLI (c) 2025 Sayantan Santra\n")
		fmt.Fprintf(writer, "Command line arguments to override config file values:\n")

		flag.VisitAll(func(f *flag.Flag) {
			fmt.Fprintf(writer, "\t--%v\t%v\n", f.Name, f.Usage)
		})
		writer.Flush()
	}

	flag.Parse()
	if *urlFlag != "" {
		config.URL = *urlFlag
	}
	if *apiFlag != "" {
		config.APIKey = *apiFlag
	}

	return config
}
