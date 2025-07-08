package internal

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"text/tabwriter"
)

const version = "0.4.1"

type Config struct {
	URL      string `json:"url"`
	APIKey   string `json:"apiKey"`
	Password string `json:"password"`
}

type AppData struct {
	Config     Config
	Subcommand string
	Input1     string
	Input2     string
	Input3     string
}

func ParseData() AppData {
	log.SetFlags(0)

	urlFlag := flag.String("url", "", "URL of the Chhoto URL server.")
	apiFlag := flag.String("api-key", "", "API Key of the Chhoto URL server (preferred).")
	passFlag := flag.String("password", "", "Password for the Chhoto URL server. It may also be passed interactively.")
	flag.BoolFunc("version", "Prints the version.", func(_ string) error {
		fmt.Print("v", version, "\n")
		os.Exit(0)
		return nil
	})
	flag.Usage = func() {
		writer := tabwriter.NewWriter(flag.CommandLine.Output(), 0, 4, 4, ' ', 0)
		fmt.Fprintf(writer, "Chhoto URL CLI (c) 2025 Sayantan Santra\n")
		fmt.Fprint(writer, "By default, config will be loaded from $XDG_CONFIG_HOME/chhoto/config.json\n")
		fmt.Fprint(writer, "But these can be overridden by using the flags.\n")

		fmt.Fprint(writer, "Subcommands:\n")
		fmt.Fprint(writer, "\tnew <longurl> [<shorturl>] [<expiry-delay>]\tCreate a new shorturl.\n")
		fmt.Fprint(writer, "\t\tIf shorturl is not provided, it will be generated automatically.\n")
		fmt.Fprint(writer, "\t\tExpiry delay should be in seconds. Default value is 0, which means no expiry.\n")
		fmt.Fprint(writer, "\tdelete <shorturl>\tDelete a given shorturl.\n")
		fmt.Fprint(writer, "\texpand <shorturl>\tGet info about a particular shorturl.\n")
		fmt.Fprint(writer, "\tgetall\tGet info about all shorturls in the server.\n")
		fmt.Fprint(writer, "\tgetconfig\tPrint the backend config.\n")

		fmt.Fprint(writer, "Flags:\n")
		flag.VisitAll(func(f *flag.Flag) {
			fmt.Fprintf(writer, "\t--%v\t%v\n", f.Name, f.Usage)
		})
		writer.Flush()
	}
	err := parseFlags()
	if err != nil {
		log.Fatalln("There was an error parsing command line flags.")
	}

	config := Config{}
	if *urlFlag == "" {
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
	} else {
		// Since an url was provided, we assume that new details
		// are needed which are not in config
		config.URL = *urlFlag
		config.APIKey = *apiFlag
		config.Password = *passFlag
		if *apiFlag == "" && *passFlag == "" {
			var pass string
			fmt.Println("No API key was password was provided.")
			fmt.Print("Type your password here: ")
			pass = readPass()
			config.Password = pass
		}
	}

	args := flag.Args()
	var arg1, arg2, arg3, arg4 string
	if len(args) == 0 {
		log.Fatalln("No subcommand was supplied! Please see help.")
	}
	if len(args) >= 1 {
		arg1 = args[0]
		arg2 = ""
		arg3 = ""
		arg4 = "0"
	}
	if len(args) >= 2 {
		arg2 = args[1]
	}
	if len(args) >= 3 {
		arg3 = args[2]
	}
	if len(args) >= 4 {
		arg4 = args[3]
	}
	if len(args) >= 5 {
		log.Fatalln("Too many arguments were supplied! Please see help.")
	}

	return AppData{
		Config:     config,
		Subcommand: arg1,
		Input1:     arg2,
		Input2:     arg3,
		Input3:     arg4,
	}
}

func parseFlags() error {
	return parseFlagSet(flag.CommandLine, os.Args[1:])
}

// ParseFlagSet works like flagset.Parse(), except positional
// args and flag args can be specified in any order.
func parseFlagSet(flagset *flag.FlagSet, args []string) error {
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
