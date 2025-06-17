package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"text/tabwriter"
)

const version = "0.1.0"

type Config struct {
	URL    string `json:"url"`
	APIKey string `json:"apiKey"`
}

type AppData struct {
	Config     Config
	Subcommand string
	Input1     string
	Input2     string
}

func ParseData() AppData {
	log.SetFlags(0)

	urlFlag := flag.String("url", "", "URL of the Chhoto URL server.")
	apiFlag := flag.String("api-key", "", "API Key of the Chhoto URL server.")
	flag.BoolFunc("version", "Prints the version.", func(_ string) error {
		fmt.Print("v", version, "\n")
		os.Exit(0)
		return nil
	})
	flag.Usage = func() {
		writer := tabwriter.NewWriter(flag.CommandLine.Output(), 0, 4, 4, ' ', 0)
		fmt.Fprintf(writer, "Chhoto URL CLI (c) 2025 Sayantan Santra\n")
		fmt.Fprintf(writer, "Command line arguments to override config file values:\n")

		flag.VisitAll(func(f *flag.Flag) {
			fmt.Fprintf(writer, "\t--%v\t%v\n", f.Name, f.Usage)
		})
		writer.Flush()
	}
	err := ParseFlags()
	if err != nil {
		log.Fatalln("There was an error parsing command line flags.")
	}

	config := Config{}
	if *urlFlag == "" || *apiFlag == "" {
		configDir, ok := os.LookupEnv("XDG_CONFIG_HOME")
		if !ok {
			configDir = "~/.config"
		}
		file, err := os.Open(configDir + "/chhoto/config.json")
		if err != nil {
			log.Fatalln("Could not load config from " + configDir + "/chhoto/config.json. Quitting!")
		}
		defer file.Close()
		decoder := json.NewDecoder(file)
		err = decoder.Decode(&config)
		if err != nil {
			fmt.Println("error:", err)
		}
	}

	if *urlFlag != "" {
		config.URL = *urlFlag
	}
	if *apiFlag != "" {
		config.APIKey = *apiFlag
	}

	args := flag.Args()
	var arg1, arg2, arg3 string
	if len(args) == 0 {
		log.Fatalln("No subcommand was supplied! Please see help.")
	}
	if len(args) >= 1 {
		arg1 = args[0]
		arg2 = ""
		arg3 = ""
	}
	if len(args) >= 2 {
		arg2 = args[1]
	}
	if len(args) >= 3 {
		arg3 = args[2]
	}
	if len(args) >= 4 {
		log.Fatalln("Too many arguments were supplied! Please see help.")
	}

	return AppData{
		Config:     config,
		Subcommand: arg1,
		Input1:     arg2,
		Input2:     arg3,
	}
}

func ParseFlags() error {
	return ParseFlagSet(flag.CommandLine, os.Args[1:])
}

// ParseFlagSet works like flagset.Parse(), except positional
// args and flag args can be specified in any order.
func ParseFlagSet(flagset *flag.FlagSet, args []string) error {
	var positionalArgs []string
	for {
		if err := flagset.Parse(args); err != nil {
			return err
		}
		// Consume all the flags that were parsed as flags.
		args = args[len(args)-flagset.NArg():]
		if len(args) == 0 {
			break
		}
		positionalArgs = append(positionalArgs, args[0])
		args = args[1:]
	}
	return flagset.Parse(positionalArgs)
}
