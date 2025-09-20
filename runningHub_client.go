package runninghub_tools

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"time"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gclient"
)

var (
	statusOK = 200
)

type Client struct {
	ApiKey string `json:"api_key"`
}

func NewClient(apiKey string) *Client {
	return &Client{
		ApiKey: apiKey,
	}
}

func (r *Client) doPost(url string, reqBody []byte) (resp *gclient.Response, err error) {
	resp, err = g.Client().
		SetTimeout(60*time.Second).
		Header(map[string]string{"Content-Type": "application/json", "Host": "www.runninghub.cn"}).
		Post(context.Background(), url, reqBody)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// 获取api账户信息
func (r *Client) GetAccountStatus() (res *RunningHubAccountSuccessResponseData, err error) {
	res = &RunningHubAccountSuccessResponseData{
		CurrentTaskCounts: "0",
		RemainCoins:       "0",
	}
	url := "https://www.runninghub.cn/uc/openapi/accountStatus"
	payloadData := map[string]string{
		"apikey": r.ApiKey,
	}

	reqBody, _ := json.Marshal(payloadData)

	resp, err := r.doPost(url, reqBody)
	if err != nil {
		return nil, err
	}
	defer resp.Close()

	// 检查响应状态码
	if resp.StatusCode != statusOK {
		return nil, gerror.Newf("unexpected status code: %d", resp.StatusCode)
	}
	var response *RunningHubResponse[json.RawMessage]
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	if !response.OK() {
		responseJson, _ := json.Marshal(response)
		err = gerror.NewCodef(gcode.New(-1, "获取任务生成结果时发生错误", response), "请求的响应码不符合要求（应为0，实际为%d）, 错误详情: %s", response.Code, responseJson)
		return nil, err
	}
	if err := json.Unmarshal(response.Data, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (r *Client) CreateTask(payloadData *CreateTaskRequestInfo) (res *RunningHubCreateTaskSuccessResponseData, err error) {
	res = &RunningHubCreateTaskSuccessResponseData{}
	if payloadData == nil {
		return nil, gerror.New("Running Create Task Request`s Payload data cannot be empty")
	}
	url := "https://www.runninghub.cn/task/openapi/create"
	payloadData.ApiKey = r.ApiKey
	reqBody, _ := json.Marshal(payloadData)
	resp, err := r.doPost(url, reqBody)
	if err != nil {
		return nil, err
	}
	defer resp.Close()

	// 检查响应状态码
	if resp.StatusCode != statusOK {
		return nil, gerror.Newf("unexpected status code: %d", resp.StatusCode)
	}

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var response *RunningHubCreateTaskSuccessResponse
	json.Unmarshal(body, &response)
	if response.Code != 0 {
		responseJson, _ := json.Marshal(response)
		err = gerror.NewCodef(gcode.New(-1, "获取任务生成结果时发生错误", response), "请求的响应码不符合要求（应为0，实际为%d）, 错误详情: %s", response.Code, responseJson)
		return nil, err
	}
	return response.Data, nil
}

// 获取任务状态
func (r *Client) GetTaskStatus(taskId string) (res string) {
	if taskId == "" {
		log.Fatal("task_id cannot be empty")
	}
	url := "https://www.runninghub.cn/task/openapi/status"

	// 准备请求数据
	payloadData := map[string]string{
		"taskId": taskId,
		"apiKey": r.ApiKey,
	}

	reqBody, _ := json.Marshal(payloadData)
	resp, err := r.doPost(url, reqBody)
	if err != nil {
		return ""
	}
	defer resp.Close()

	// 检查响应状态码
	if resp.StatusCode != statusOK {
		return ""
	}

	var response *RunningHubResponse[json.RawMessage]
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return ""
	}
	if err := json.Unmarshal(response.Data, &res); err != nil {
		return ""
	}
	return res
}

// 获取任务结果
func (r *Client) GetTaskResult(taskId string) (res []*RunningHubGetTaskResultSuccessResponseData, err error) {
	res = []*RunningHubGetTaskResultSuccessResponseData{}
	if taskId == "" {
		log.Fatal("task_id cannot be empty")
	}
	url := "https://www.runninghub.cn/task/openapi/outputs"
	// 准备请求数据
	payloadData := map[string]string{
		"taskId": taskId,
		"apiKey": r.ApiKey,
	}

	reqBody, _ := json.Marshal(payloadData)
	resp, err := r.doPost(url, reqBody)
	if err != nil {
		return nil, err
	}
	defer resp.Close()

	// 检查响应状态码
	if resp.StatusCode != statusOK {
		return nil, gerror.Newf("unexpected status code: %d", resp.StatusCode)
	}
	var response *RunningHubResponse[json.RawMessage]
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	if !response.OK() {
		responseJson, _ := json.Marshal(response)
		err = gerror.NewCodef(gcode.New(-1, "获取任务生成结果时发生错误", response), "请求的响应码不符合要求（应为0，实际为%d）, 错误详情: %s", response.Code, responseJson)
		return nil, err
	}
	if err := json.Unmarshal(response.Data, &res); err != nil {
		return nil, err
	}
	return res, nil
}

// 取消任务
func (r *Client) CancelTask(taskId string) (err error) {
	if taskId == "" {
		log.Fatal("task_id cannot be empty")
	}

	url := "https://www.runninghub.cn/task/openapi/cancel"

	// 准备请求数据
	payloadData := map[string]string{
		"taskId": taskId,
		"apiKey": r.ApiKey,
	}

	reqBody, _ := json.Marshal(payloadData)
	resp, err := r.doPost(url, reqBody)
	if err != nil {
		return err
	}
	defer resp.Close()

	// 检查响应状态码
	if resp.StatusCode != statusOK {
		return gerror.Newf("unexpected status code: %d", resp.StatusCode)
	}
	var response *RunningHubResponse[json.RawMessage]
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return err
	}
	if !response.OK() {
		responseJson, _ := json.Marshal(response)
		err = gerror.NewCodef(gcode.New(-1, "获取任务生成结果时发生错误", response), "请求的响应码不符合要求（应为0，实际为%d）, 错误详情: %s", response.Code, responseJson)
		return err
	}
	return nil
}
