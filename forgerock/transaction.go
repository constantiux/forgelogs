package forgerock

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/tidwall/gjson"
)

func FetchException(Config ConfigData, TransactionId string, EndTime string) string {
	resource := "/monitoring/logs"
	params := url.Values{}
	params.Add("source", Config.FrSource)
	params.Add("endTime", EndTime)
	params.Add("transactionId", TransactionId)
	params.Add("_queryFilter", "/payload/exception pr")

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
	resp, err := client.Do(req)
	CheckError(err)

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("An error occured! Received %v status code from ForgeRock.", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	CheckError(err)
	//fmt.Println(string(body))

	var payload ForgeRockLogModel
	err = json.Unmarshal(body, &payload)
	body = nil
	CheckError(err)

	Result, err := json.Marshal(payload.Result[0])
	CheckError(err)
	Payload := gjson.Get(string(Result), "payload")
	ExceptionMessage := Payload.Get("exception")

	if ExceptionMessage.Exists() {
		ExceptionDetail := strings.Split(ExceptionMessage.Str, "\n")
		Reason := ExceptionDetail[0]

		return Reason
	} else {
		return ""
	}
}
