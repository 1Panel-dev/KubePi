package prometheus

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Prometheus struct {
	Address  string `json:"address"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewPrometheusClient(address, username, password string) *Prometheus {
	return &Prometheus{
		Address:  address,
		Username: username,
		Password: password,
	}
}

func (p *Prometheus) GetMetrics(address, username, password string) (data *MetricsExplorerResp, err error) {
	reqUrl := address + "/api/v1/label/__name__/values"
	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(username, password)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New("获取Prometheus Metrics列表失败,当前状态码为: " + resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("获取Prometheus Metrics列表失败: " + err.Error())
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, errors.New("解析获取Prometheus Metrics列表json数据失败: " + err.Error())
	}

	return data, nil
}

func (p *Prometheus) TestConnect(address, username, password string) error {
	req, err := http.NewRequest("GET", address, nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(username, password)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New("请求Prometheus地址失败,当前状态码为: " + resp.Status)
	}
	return nil
}

func (p *Prometheus) QueryMetrics(address, username, password, promql, timestamp string) (data []*ProcessedMetric, err error) {
	params := url.Values{}
	params.Set("query", promql)
	timestamp1, _ := strconv.ParseFloat(timestamp, 64)
	params.Set("time", fmt.Sprintf("%.3f", timestamp1))
	reqUrl := fmt.Sprintf("%s/api/v1/query?%s", address, params.Encode())
	fmt.Println("reqUrl：", reqUrl)
	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(username, password)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New("执行Prometheus Promql查询失败,当前状态码为: " + resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("执行Prometheus Promql查询失败: " + err.Error())
	}

	var queryMetricsResp QueryMetricsResp
	err = json.Unmarshal(body, &queryMetricsResp)
	if err != nil {
		return nil, errors.New("解析Prometheus Promql json数据失败: " + err.Error())
	}

	var processedMetrics []*ProcessedMetric

	for _, result := range queryMetricsResp.Data.Result {
		labels := []string{}
		metricName := ""

		for key, value := range result.Metric {
			if key == "__name__" {
				metricName = value
			} else {
				labels = append(labels, fmt.Sprintf(`%s="%s"`, key, value))
			}
		}

		labelString := strings.Join(labels, ", ")
		promql = fmt.Sprintf(`%s{%s}`, metricName, labelString)
		metricValue, _ := strconv.ParseFloat(result.Value[1].(string), 64)

		processedMetrics = append(processedMetrics, &ProcessedMetric{
			Metrics: promql,
			Value:   metricValue,
		})
	}
	return processedMetrics, nil
}
