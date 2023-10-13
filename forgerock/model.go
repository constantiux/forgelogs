package forgerock

import (
	"time"
)

type ConfigData struct {
	TenantUrl       *string   `json:"tenantUrl,omitempty"`
	ApiKey          *string   `json:"apiKey,omitempty"`
	ApiSecret       *string   `json:"apiSecret,omitempty"`
	FrTrees         *[]string `json:"trees"`
	FrSource        string
	BeginTime       time.Time
	EndTime         time.Time
	FrTreeName      string
	FailedOnlyFlag  bool
	TransactionFlag bool
	FilterTreesFlag bool
	AllTreesFlag    bool
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
