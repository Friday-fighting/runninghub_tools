package runninghub_tools

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gclient"
)

const (
	getAccountStatus = "/uc/openapi/accountStatus"
	getTaskStatus    = "/task/openapi/status"
	createTask       = "/task/openapi/create"
	getTaskResult    = "/task/openapi/outputs"
	cancelTask       = "/task/openapi/cancel"
)

type RunningHubClient struct {
	url        string
	ApiKey     string `json:"api_key"`
	httpClient *gclient.Client
}

func NewClient(in *RunningHubClientConfig) *RunningHubClient {
	if in.Host == "" {
		in.Host = "www.runninghub.cn"
	}
	protocol := "https"
	if in.UseHttpReq {
		protocol = "http"
	}
	url := fmt.Sprintf("%s://%s", protocol, in.Host)
	return &RunningHubClient{
		ApiKey: in.ApiKey,
		url:    url,
		httpClient: g.Client().
			SetTimeout(60*time.Second).
			SetHeader("Content-Type", "application/json").
			SetHeader("Host", in.Host),
	}
}

func (c *RunningHubClient) doPost(ctx context.Context, url string, reqBody []byte) (res *RunningHubResponse, err error) {
	httpClient := c.httpClient.Clone()
	resp, err := httpClient.Post(ctx, url, reqBody)
	if err != nil {
		return nil, err
	}
	defer resp.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return nil, gerror.Newf("unexpected status code: %d", resp.StatusCode)
	}
	var response *RunningHubResponse
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	return response, nil
}

// GetAccountStatus 获取api账户信息
func (c *RunningHubClient) GetAccountStatus(ctx context.Context) (res *GetAccountRes, err error) {
	url := fmt.Sprintf("%s%s", c.url, getAccountStatus)
	reqBody, _ := json.Marshal(g.Map{
		"apikey": c.ApiKey,
	})
	resp, err := c.doPost(ctx, url, reqBody)
	if err != nil {
		return nil, err
	}
	if resp.Code != 0 {
		return nil, gerror.Newf("GetAccountStatus fail, code: %d, msg: %s", resp.Code, resp.Msg)
	}
	if err := json.Unmarshal(resp.Data, &res); err != nil {
		return nil, fmt.Errorf("decode success data fail: %w", err)
	}
	return res, nil
}

// CreateTask 创建任务
func (c *RunningHubClient) CreateTask(ctx context.Context, payloadData *CreateTaskReq) (res *CreateTaskRes, err error) {
	if payloadData == nil || payloadData.WorkflowId == "" || len(payloadData.NodeInfoList) == 0 {
		return nil, gerror.New("Running Create Task Request`s Payload data params of WorkflowId and NodeInfoList must be not empty")
	}
	url := fmt.Sprintf("%s%s", c.url, createTask)
	payloadData.ApiKey = c.ApiKey
	reqBody, _ := json.Marshal(payloadData)

	resp, err := c.doPost(ctx, url, reqBody)
	if err != nil {
		return nil, err
	}
	if resp.Code != 0 {
		return nil, gerror.Newf("CreateTask fail, code: %d, msg: %s", resp.Code, resp.Msg)
	}
	if err := json.Unmarshal(resp.Data, &res); err != nil {
		return nil, fmt.Errorf("decode success data fail: %w", err)
	}
	return res, nil
}

// GetTaskStatus 获取任务状态
func (c *RunningHubClient) GetTaskStatus(ctx context.Context, taskId string) (res string, err error) {
	if taskId == "" {
		return "", gerror.New("task_id cannot be empty")
	}
	url := fmt.Sprintf("%s%s", c.url, getTaskStatus)
	reqBody, _ := json.Marshal(g.Map{
		"taskId": taskId,
		"apiKey": c.ApiKey,
	})
	resp, err := c.doPost(ctx, url, reqBody)
	if err != nil {
		return "", err
	}
	if resp.Code != 0 {
		return "", gerror.Newf("GetTaskStatus fail, code: %d, msg: %s", resp.Code, resp.Msg)
	}
	if err := json.Unmarshal(resp.Data, &res); err != nil {
		return "", fmt.Errorf("decode success data fail: %w", err)
	}
	return res, nil
}

// GetTaskResult 获取任务结果
func (c *RunningHubClient) GetTaskResult(ctx context.Context, taskId string) (res *GetTaskResultRes, err error) {
	if taskId == "" {
		return nil, gerror.New("task_id cannot be empty")
	}
	url := fmt.Sprintf("%s%s", c.url, getTaskResult)
	reqBody, _ := json.Marshal(g.Map{
		"taskId": taskId,
		"apiKey": c.ApiKey,
	})
	resp, err := c.doPost(ctx, url, reqBody)
	if err != nil {
		return nil, err
	}
	res = &GetTaskResultRes{
		SuccessItems: []*SuccessOfGetTaskResultResponseData{},
		FailedReason: &FailedOfGetTaskResultResponseData{},
		Code:         resp.Code,
		Msg:          resp.Msg,
	}
	if resp.Code == 0 {
		if err := json.Unmarshal(resp.Data, &res.SuccessItems); err != nil {
			return nil, fmt.Errorf("decode success data fail: %w", err)
		}
		return res, nil
	}
	var failData struct {
		FailedReason FailedOfGetTaskResultResponseData `json:"failedReason"`
	}
	if err := json.Unmarshal(resp.Data, &failData); err != nil {
		return nil, fmt.Errorf("decode fail data fail: %w", err)
	}
	res.FailedReason = &failData.FailedReason
	return res, nil
}

// 取消任务
func (c *RunningHubClient) CancelTask(ctx context.Context, taskId string) (err error) {
	if taskId == "" {
		log.Fatal("task_id cannot be empty")
	}
	url := fmt.Sprintf("%s%s", c.url, cancelTask)
	reqBody, _ := json.Marshal(g.Map{
		"taskId": taskId,
		"apiKey": c.ApiKey,
	})
	resp, err := c.doPost(ctx, url, reqBody)
	if err != nil {
		return err
	}
	if resp.Code != 0 {
		return gerror.Newf("CancelTask fail, code: %d, msg: %s", resp.Code, resp.Msg)
	}
	return nil
}

// GetTaskStatusAndResult 获取任务状态和结果, 支持获取结果时重试
func (c *RunningHubClient) GetTaskStatusAndResult(ctx context.Context, in *GetTaskStatusAndResultReq) (res *GetTaskStatusAndResultRes, err error) {
	if in.TaskId == "" {
		return nil, gerror.New("task_id cannot be empty")
	}
	status, err := c.GetTaskStatus(ctx, in.TaskId)
	if err != nil {
		return nil, err
	}
	res = &GetTaskStatusAndResultRes{
		Status: status,
	}
	switch status {
	case "SUCCESS", "FAILED":
		var result *GetTaskResultRes
		if in.MaxTries <= 0 {
			in.MaxTries = 5
		}
		if in.SleepTime <= 0 {
			in.SleepTime = 10
		} else if in.SleepTime >= 60 {
			in.SleepTime = 60
		}
		for i := 0; i < in.MaxTries; i++ {
			result, err = c.GetTaskResult(ctx, in.TaskId)
			if err == nil && result != nil {
				res.Code = result.Code
				res.Msg = result.Msg
				res.SuccessItems = result.SuccessItems
				res.FailedReason = result.FailedReason
				return res, nil
			}
			time.Sleep(time.Duration(in.SleepTime) * time.Second)
		}
	}
	return res, nil
}
