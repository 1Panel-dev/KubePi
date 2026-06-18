package grafana

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type Grafana struct {
	Address             string `json:"address"`
	Enable              bool   `json:"enable"`
	DefaultDashboard    bool   `json:"default_dashboard"`
	ServiceAccountToken string `json:"service_account_token"`
}

func NewGrafanaClient(address, token string, enable, defaultDashboard bool) *Grafana {
	return &Grafana{
		Address:             address,
		Enable:              enable,
		DefaultDashboard:    defaultDashboard,
		ServiceAccountToken: token,
	}
}

// 强制固定 "KubePi Dashboards" 文件夹的UID
var folderUid = "nErXDvCkzzkubepi"

// 测试Grafana连接是否能访问
func (g *Grafana) TestConnect(address string) error {
	req, err := http.NewRequest("GET", address, nil)
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
		return errors.New("请求Grafana地址失败,当前状态码为: " + resp.Status)
	}
	return nil
}

// 获取文件夹列表
func (g *Grafana) GetFolders(address, token string) (data []*GetFoldersResp, err error) {
	reqUrl := address + "/api/folders"
	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 判断是否增加成功
	if resp.StatusCode != 200 {
		return nil, errors.New("获取Grafana文件夹列表失败,当前状态码为: " + resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("获取Grafana文件夹列表失败: " + err.Error())
	}

	// 解析成json格式
	var respData []*GetFoldersResp
	err = json.Unmarshal(body, &respData)
	if err != nil {
		return nil, errors.New("解析Grafana文件夹列表json数据失败: " + err.Error())
	}

	return respData, nil
}

// 创建文件夹
func (g *Grafana) CreateFolder(address, token, folderName string) error {
	folders, err := g.GetFolders(address, token)
	if err != nil {
		return err
	}
	for _, v := range folders {
		if v.Uid == folderUid {
			return nil
		}
	}

	params := new(struct {
		Uid   string `json:"uid"`
		Title string `json:"title"`
	})
	params.Uid = folderUid
	params.Title = folderName
	paramsData, _ := json.Marshal(params)

	reqUrl := address + "/api/folders"
	req, err := http.NewRequest("POST", reqUrl, bytes.NewBuffer(paramsData))
	if err != nil {
		return errors.New("创建Grafana文件夹失败: " + err.Error())
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.New("创建Grafana文件夹失败: " + err.Error())
	}
	defer resp.Body.Close()

	// 判断是否增加成功
	if resp.StatusCode != 200 {
		return errors.New("创建Grafana文件夹失败,当前状态码为: " + resp.Status)
	}
	return nil
}

// 获取仪表盘
func (g *Grafana) GetDashboards(address, token, folderUid string) (data []*GetDashboardsResp, err error) {
	reqUrl := address + "/api/search?limit=50&type=dash-db&folderUid=" + folderUid
	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 判断是否增加成功
	if resp.StatusCode != 200 {
		return nil, errors.New("获取Grafana仪表盘列表失败,当前状态码为: " + resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("获取Grafana仪表盘列表失败: " + err.Error())
	}

	// 解析成json格式
	var respData []*GetDashboardsResp
	err = json.Unmarshal(body, &respData)
	if err != nil {
		return nil, errors.New("解析Grafana仪表盘列表json数据失败: " + err.Error())
	}

	return respData, nil
}

/*
导入仪表盘
仪表盘UID自定义：
Namespace Overview：NamespaceOverviewKubePi
*/
func (g *Grafana) ImportDashboards(address, token, uid, jsonContent string) error {
	dashboards, err := g.GetDashboards(address, token, folderUid)
	if err != nil {
		return err
	}
	for _, v := range dashboards {
		if v.Uid == uid {
			return nil
		}
	}

	params := struct {
		Dashboard interface{} `json:"dashboard"`
		Overwrite bool        `json:"overwrite"`
		FolderUid string      `json:"folderUid"`
	}{
		Overwrite: true,
		FolderUid: folderUid,
	}

	// 解析文件中的JSON并将其赋值给params.Dashboard
	err = json.Unmarshal([]byte(jsonContent), &params.Dashboard)
	if err != nil {
		return errors.New("解析jsonContent仪表盘JSON数据失败: " + err.Error())
	}
	paramsData, err := json.Marshal(params)
	if err != nil {
		return errors.New("导入Grafana仪表盘,序列化请求参数失败: " + err.Error())
	}
	reqUrl := address + "/api/dashboards/import"
	req, err := http.NewRequest("POST", reqUrl, bytes.NewBuffer(paramsData))
	if err != nil {
		return errors.New("导入Grafana仪表盘失败: " + err.Error())
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.New("导入Grafana仪表盘失败: " + err.Error())
	}
	defer resp.Body.Close()

	// 判断是否增加成功
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return errors.New("导入Grafana仪表盘失败, 当前状态码为: " + resp.Status + ", 错误信息: " + string(body))
	}
	return nil
}
