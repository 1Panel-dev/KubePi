package prometheus

type MetricsExplorerResp struct {
	Status string   `json:"status"`
	Data   []string `json:"data"`
}

type QueryMetricsResp struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric map[string]string `json:"metric"`
			Value  []interface{}     `json:"value"`
		} `json:"result"`
	} `json:"data"`
}

// 定义新的结构体用于存储处理后的数据
type ProcessedMetric struct {
	Metrics string  `json:"metrics"`
	Value   float64 `json:"value"`
}
