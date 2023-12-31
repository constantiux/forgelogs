package forgerock

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/tidwall/gjson"
)

func stringSliceContains(elems []string, v string) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}

func GenerateReport(Config ConfigData) {
	var OutputCSV = [][]string{
		{"id", "timestamp", "treeName", "transactionId", "actor", "outcome", "reason"},
	}
	var CSVIndex int = -1

	for i := 0; i < OutputLog.ResultCount; i++ {
		SubResult, err := json.Marshal(OutputLog.Result[i])
		CheckError(err)
		Payload := gjson.Get(string(SubResult), "payload")
		if Payload.Exists() {
			PayloadTreeName := Payload.Get("entries.0.info.treeName").String()
			if (len(Config.FrTreeName) == 0) && Config.FilterTreesFlag && !stringSliceContains(*Config.FrTrees, PayloadTreeName) {
				continue
			}

			PayloadTimestamp := Payload.Get("timestamp").String()
			PayloadTransactionId := Payload.Get("transactionId").String()
			PayloadActor := Payload.Get("principal.0").String()
			PayloadOutcome := Payload.Get("result").String()
			var PayloadReason string = ""

			if (strings.Compare(PayloadOutcome, "FAILED") == 0) && (Config.TransactionFlag) {
				fmt.Printf(".")
				t, err := time.Parse(time.RFC3339, PayloadTimestamp)
				if err != nil {
					CheckError(err)
				}
				PayloadReason = FetchException(Config, PayloadTransactionId, t.Add(time.Hour*1).Format(time.RFC3339))
			}

			CSVIndex += 1
			NewEntry := []string{strconv.Itoa(CSVIndex), PayloadTimestamp, PayloadTreeName, PayloadTransactionId, PayloadActor, PayloadOutcome, PayloadReason}
			OutputCSV = append(append([][]string{}, OutputCSV...), NewEntry)
		}
	}

	if CSVIndex > -1 {
		fmt.Println()
		log.Printf("Generating a log report with applied filter (if any)..\n\n")

		cwd, _ := os.Getwd()
		LogFileTimestamp = strconv.Itoa(int(time.Now().Unix()))
		FullFileName := fmt.Sprintf("%s_%s", OutFileName, LogFileTimestamp)
		f, err := os.Create(fmt.Sprintf("%s/%s.csv", cwd, FullFileName))
		CheckError(err)
		writer := csv.NewWriter(f)
		err = writer.WriteAll(OutputCSV)
		CheckError(err)

		log.Printf("Total report entries: %v\n", CSVIndex+1)
		log.Printf("Output file %v has been successfully generated.\n", FullFileName+".csv")
	}
}
