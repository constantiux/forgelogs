package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"constantiux/forgelogs/forgerock"
	"constantiux/forgelogs/myparser"

	"github.com/akamensky/argparse"
)

var Config forgerock.ConfigData
var InFileName string = "config.json"

func CheckError(e error) {
	if e != nil {
		log.Fatalln(e)
	}
}

func LoadConfigFile() {
	cwd, _ := os.Getwd()
	content, err := os.ReadFile(cwd + "/" + InFileName)
	CheckError(err)
	err = json.Unmarshal(content, &Config)
	CheckError(err)
	if Config.TenantUrl == nil || Config.ApiKey == nil || Config.ApiSecret == nil {
		log.Fatalln("Config file is corrupted, should contain tenantUrl, apiKey, and apiSecret")
	}
	// Check if AllTreesFlag is on and required config is present
	if Config.FrTrees == nil && Config.FilterTreesFlag {
		log.Fatalln("Config file is missing \"trees\", required as --filter-trees is specified")
	}
}

func LoadForgerockConfig(FrTreeName string, FailedOnlyFlag bool, TransactionFlag bool, FilterTreesFlag bool, AllTreesFlag bool) {
	Config.FrSource = myparser.FrSource
	Config.BeginTime = myparser.BeginTime
	Config.EndTime = myparser.EndTime
	Config.FrTreeName = FrTreeName
	Config.FailedOnlyFlag = FailedOnlyFlag
	Config.TransactionFlag = TransactionFlag
	Config.FilterTreesFlag = FilterTreesFlag
	Config.AllTreesFlag = AllTreesFlag
}

func main() {
	// Create new parser object
	parser := argparse.NewParser("forgelogs", "Introducing Forgelogs: your ForgeRock logging companion.\n All rights reserved to the original author.")

	// Create flags
	_ = parser.Selector("s", "source", []string{"am", "idm"}, &argparse.Options{Required: true, Validate: myparser.ValidateArgSource, Help: "Accepts am or idm"})

	_ = parser.String("b", "begin", &argparse.Options{Required: true, Validate: myparser.ValidateArgBeginTime,
		Help: "Start datetime with RFC3339 format, examples:\n" + myparser.GenerateArgTimeExamples(false)})

	_ = parser.String("e", "end", &argparse.Options{Required: false, Validate: myparser.ValidateArgEndTime,
		Help: "(Optional, alt. to -d/--duration) End datetime with RFC3339 format, examples:\n" + myparser.GenerateArgTimeExamples(true)})

	_ = parser.String("d", "duration", &argparse.Options{Required: false, Validate: myparser.ValidateArgDuration,
		Help: "(Optional, alt. to -e/--end) Time elapsed since the start time using N(s|m|h|d), e.g. 7d for 7 days"})

	var FrTreeName *string = parser.String("t", "tree", &argparse.Options{Required: false, Help: "Filter logs using Tree/journey identifier"})

	var FailedOnlyFlag *bool = parser.Flag("", "failed-only", &argparse.Options{Required: false, Help: "Specify if only require failed results"})

	var TransactionFlag *bool = parser.Flag("", "detailed", &argparse.Options{Required: false, Help: "Specify if only require details on failed results"})

	var FilterTreesFlag *bool = parser.Flag("", "filter-trees", &argparse.Options{Required: false, Help: "Specify if require logging on selected journeys"})

	var AllTreesFlag *bool = parser.Flag("", "all-trees", &argparse.Options{Required: false, Help: "Specify if require logging on all journeys"})

	// Parse input
	err := parser.Parse(os.Args)
	if err != nil {
		// In case of error print error and print usage
		// This can also be done by passing -h or --help flags
		fmt.Printf("%v\n", err)
		//fmt.Print(parser.Usage(err))
		os.Exit(0)
	}

	// Unload parsed user inputs to config
	LoadForgerockConfig(*FrTreeName, *FailedOnlyFlag, *TransactionFlag, *FilterTreesFlag, *AllTreesFlag)

	// Check if user-defined config file exists
	LoadConfigFile()

	if *FrTreeName != "" && (*FilterTreesFlag || *AllTreesFlag) {
		log.Println("Reminder: the individual tree specified is overwriting the other option(s)")
	} else if *FilterTreesFlag && *AllTreesFlag {
		log.Println("Reminder: --filter-trees is overwriting --all-trees")
	}

	if *FailedOnlyFlag && *FrTreeName == "" && !*FilterTreesFlag && !*AllTreesFlag {
		log.Println("Reminder: --failed-only is ignored as no tree is specified")
	}

	// Initiate script
	log.Println("Tenant:", *Config.TenantUrl)
	log.Println("Source:", Config.FrSource)
	log.Println("Start at:", Config.BeginTime.Format(time.RFC3339))
	log.Println("End at:", Config.EndTime.Format(time.RFC3339))
	fmt.Println("")
	log.Println("Please wait, currently generating logs...")
	log.Println("Do not terminate while the program is running.")
	fmt.Println("")

	forgerock.GenerateLogs(Config)
}
