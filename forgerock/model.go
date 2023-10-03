package forgerock

import (
	"time"
)

type ConfigData struct {
	TenantUrl      *string `json:"tenantUrl,omitempty"`
	ApiKey         *string `json:"apiKey,omitempty"`
	ApiSecret      *string `json:"apiSecret,omitempty"`
	FrSource       string
	BeginTime      time.Time
	EndTime        time.Time
	FrTreeName     string
	FailedOnlyFlag bool
}

type ForgeRockLogModel struct {
	Result                  []map[string]interface{} `json:"result"`
	ResultCount             int                      `json:"resultCount"`
	PagedResultsCookie      string                   `json:"pagedResultsCookie"`
	TotalPagedResultsPolicy string                   `json:"totalPagedResultsPolicy"`
	TotalPagedResults       int                      `json:"totalPagedResults"`
	RemainingPagedResults   int                      `json:"remainingPagedResults"`
}

type OutputLogModel struct {
	Result       []map[string]interface{} `json:"result"`
	ResultCount  int                      `json:"resultCount"`
	GenerateTime string                   `json:"generatedAt"`
}
