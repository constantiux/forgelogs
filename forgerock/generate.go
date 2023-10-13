package forgerock

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

var OutputLog OutputLogModel
var OutFileName string = "output"
var GenerateLogsCookieCounter int = 0
var LogRequestFullDuration float64
var OriginalEndTime time.Time
var CurrentLastFetchedTimestamp time.Time
var LogFileTimestamp string

func CheckError(e error) {
	if e != nil {
		log.Fatalln(e)
	}
}

func CheckGenerateLogProgress() {
	if OriginalEndTime.IsZero() || CurrentLastFetchedTimestamp.IsZero() {
		log.Println("Still on it.. parsing more and more logs")
	} else {
		RemainingDuration := OriginalEndTime.Sub(CurrentLastFetchedTimestamp).Seconds()
		CurrentPercentage := (1 - (RemainingDuration / LogRequestFullDuration)) * 100
		log.Printf("Fetch progress: %.0f%%\n", CurrentPercentage)
	}
}

func GenerateQueryFilter(Config ConfigData) string {
	if (len(Config.FrTreeName) == 0) && !Config.FilterTreesFlag && !Config.AllTreesFlag {
		return ""
	} else {
		var TreeFilter string = "/payload/entries/info/treeName pr"
		var ResultFilter string = "/payload/result pr"

		if len(Config.FrTreeName) > 0 {
			TreeFilter = fmt.Sprintf("/payload/entries/info/treeName eq \"%s\"", Config.FrTreeName)
		}

		if Config.FailedOnlyFlag {
			ResultFilter = fmt.Sprintf("/payload/result eq \"%s\"", "FAILED")
		}

		return fmt.Sprintf("%s and %s", TreeFilter, ResultFilter)
	}
}

func GenerateLogsHelper(Config ConfigData, PagedResultsCookie string) {
	resource := "/monitoring/logs"
	params := url.Values{}
	params.Add("source", Config.FrSource)
	params.Add("beginTime", Config.BeginTime.Format(time.RFC3339))
	params.Add("endTime", Config.EndTime.Format(time.RFC3339))

	QueryFilterParam := GenerateQueryFilter(Config)
	if QueryFilterParam != "" {
		params.Add("_queryFilter", QueryFilterParam)
	}

	if PagedResultsCookie != "" {
		params.Add("_pagedResultsCookie", PagedResultsCookie)
	}

	u, _ := url.ParseRequestURI(*Config.TenantUrl)
	u.Path = resource
	u.RawQuery = params.Encode()
	urlStr := fmt.Sprintf("%v", u)

	req, err := http.NewRequest("GET", urlStr, nil)
	CheckError(err)

	req.Header.Set("Accept", "application/json")
	req.Header.Set("x-api-key", *Config.ApiKey)
	req.Header.Set("x-api-secret", *Config.ApiSecret)

	client := &http.Client{}
	fmt.Printf(".")
	resp, err := client.Do(req)
	CheckError(err)

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	CheckError(err)
	//fmt.Println(string(body))

	if resp.StatusCode != http.StatusOK {
		log.Printf("An error occured! Received %v status code from ForgeRock.", resp.StatusCode)
		log.Println(string(body))
		os.Exit(1)
	}

	var payload ForgeRockLogModel
	err = json.Unmarshal(body, &payload)
	body = nil
	CheckError(err)
	OutputLog.Result = append(append([]map[string]interface{}{}, OutputLog.Result...), payload.Result...)

	if len(payload.Result) > 0 {
		if payload.Result[len(payload.Result)-1]["timestamp"] != "" {
			CurrentLastFetchedTimestamp, _ = time.Parse(time.RFC3339, fmt.Sprint(payload.Result[len(payload.Result)-1]["timestamp"]))
		}
	}

	if GenerateLogsCookieCounter%3 == 0 {
		fmt.Println()
		CheckGenerateLogProgress()
	}

	if payload.PagedResultsCookie != "" {
		GenerateLogsCookieCounter += 1
		GenerateLogsHelper(Config, payload.PagedResultsCookie)
	} else {
		//CheckGenerateLogProgress()
		return
	}
}

func GenerateLogs(Config ConfigData) {
	TimeStarted := time.Now()
	Diff := Config.EndTime.Sub(Config.BeginTime)
	LogRequestFullDuration = Diff.Seconds()
	OriginalEndTime = Config.EndTime
	if Diff.Hours() <= 24 {
		log.Printf("Collecting logs from %v to %v\n", Config.BeginTime.Format(time.RFC3339), Config.EndTime.Format(time.RFC3339))
		GenerateLogsHelper(Config, "")
	} else {
		NDiff := int(math.Floor(Diff.Hours() / 24))
		for i := 0; i < NDiff; i++ {
			Config.EndTime = Config.BeginTime.Add(time.Hour * 24)
			if i > 0 {
				Config.EndTime = Config.EndTime.Add(-time.Nanosecond * 1)
			}
			log.Printf("Segment #%v: Collecting logs from %v to %v\n", i+1, Config.BeginTime.Format(time.RFC3339), Config.EndTime.Format(time.RFC3339))
			GenerateLogsHelper(Config, "")
			Config.BeginTime = Config.EndTime.Add(time.Nanosecond * 1)
		}
		if Config.BeginTime.Before(OriginalEndTime) {
			Config.EndTime = OriginalEndTime
			log.Printf("Segment #%v: Collecting logs from %v to %v\n", NDiff+1, Config.BeginTime.Format(time.RFC3339), Config.EndTime.Format(time.RFC3339))
			GenerateLogsHelper(Config, "")
		}
	}
	OutputLog.ResultCount = len(OutputLog.Result)

	fmt.Println()
	log.Printf("Analysing logs..\n")

	if (len(Config.FrTreeName) != 0) || Config.FilterTreesFlag || Config.AllTreesFlag {
		GenerateReport(Config)
	}

	TimeFinished := time.Now()
	OutputLog.GenerateTime = TimeFinished.UTC().Format(time.RFC3339)
	TotalTimeElapsed := TimeFinished.Sub(TimeStarted)

	body, err := json.MarshalIndent(OutputLog, "", "  ")
	CheckError(err)

	if len(LogFileTimestamp) == 0 {
		LogFileTimestamp = strconv.Itoa(int(time.Now().Unix()))
	}
	FullFileName := OutFileName + "_" + LogFileTimestamp
	if err := os.WriteFile(FullFileName+".json", body, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	fmt.Println("")
	log.Printf("Total raw log entries: %v\n", OutputLog.ResultCount)
	log.Printf("Output file %v has been successfully generated.\n", FullFileName+".json")
	log.Printf("Total time elapsed: %v\n", TotalTimeElapsed)
}
