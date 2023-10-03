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
			PayloadTimestamp := Payload.Get("timestamp").String()
			PayloadTreeName := Payload.Get("entries.0.info.treeName").String()
			PayloadTransactionId := Payload.Get("transactionId").String()
			PayloadActor := Payload.Get("principal.0").String()
			PayloadOutcome := Payload.Get("result").String()
			var PayloadReason string = ""

			if strings.Compare(PayloadOutcome, "FAILED") == 0 {
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
		log.Printf("Generating a log report..\n")

		cwd, _ := os.Getwd()
		FullFileName := fmt.Sprintf("%s_%s", OutFileName, strconv.Itoa(int(time.Now().Unix())))
		f, err := os.Create(fmt.Sprintf("%s/%s.csv", cwd, FullFileName))
		CheckError(err)
		writer := csv.NewWriter(f)
		err = writer.WriteAll(OutputCSV)
		CheckError(err)

		log.Printf("Total report entries: %v\n", CSVIndex+1)
		log.Printf("Output file %v has been successfully generated.\n", FullFileName+".csv")
	}
}
