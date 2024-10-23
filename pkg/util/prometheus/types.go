package prometheus

type MetricsExplorerResp struct {
	Status string   `json:"status"`
	Data   []string `json:"data"`
}

type QueryMetricsResp struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string        `json:"resultType"`
		Result     []interface{} `json:"result"`
	} `json:"data"`
}
