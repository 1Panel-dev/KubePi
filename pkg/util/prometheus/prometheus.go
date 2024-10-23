package prometheus

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
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

// 获取Metrics Explorer
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

	// 判断是否增加成功
	if resp.StatusCode != 200 {
		return nil, errors.New("获取Prometheus Metrics列表失败,当前状态码为: " + resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("获取Prometheus Metrics列表失败: " + err.Error())
	}

	// 解析成json格式
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, errors.New("解析获取Prometheus Metrics列表json数据失败: " + err.Error())
	}

	return data, nil
}

// 测试Prometheus连接是否能访问
func (p *Prometheus) TestConnect(address, username, password string) error {
	req, err := http.NewRequest("GET", address, nil)
	req.SetBasicAuth(username, password)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 判断是否增加成功
	if resp.StatusCode != 200 {
		return errors.New("请求Prometheus地址失败,当前状态码为: " + resp.Status)
	}
	return nil
}

// 查询
func (p *Prometheus) QueryMetrics(address, username, password, promql, timestamp string) (data *QueryMetricsResp, err error) {
	// 构建带有查询参数的 URL
	params := url.Values{}
	params.Set("query", promql)
	// 将字符串形式的 time 转换为 float64
	timestamp1, _ := strconv.ParseFloat(timestamp, 64)
	params.Set("time", fmt.Sprintf("%.3f", timestamp1)) // 时间戳保留小数点后3位
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

	// 判断是否增加成功
	if resp.StatusCode != 200 {
		return nil, errors.New("执行Prometheus Promql查询失败,当前状态码为: " + resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("执行Prometheus Promql查询失败: " + err.Error())
	}

	// 解析成json格式
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, errors.New("解析Prometheus Promql json数据失败: " + err.Error())
	}

	return data, nil
}
