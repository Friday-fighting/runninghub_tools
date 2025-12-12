package runninghub_client_utils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Friday-fighting/runninghub_tools/utility"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gclient"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/util/gconv"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	getAccountInfo  = "/uc/openapi/accountStatus"
	getTaskStatus   = "/task/openapi/status"
	createTask      = "/task/openapi/create"
	getTaskResult   = "/task/openapi/outputs"
	cancelTask      = "/task/openapi/cancel"
	getWorkflowJSON = "/api/openapi/getJsonApiFormat"
	uploadResource  = "/task/openapi/upload"
)

type RunningHubClient struct {
	url        string
	ApiKey     string `json:"api_key"`
	httpClient *gclient.Client
	Timeout    time.Duration
}

func NewClient(in *RunningHubClientConfig) *RunningHubClient {
	if in.Host == "" {
		in.Host = "www.runninghub.cn"
	}
	protocol := "https"
	if in.UseHttpReq {
		protocol = "http"
	}
	if in.Timeout <= 0 {
		in.Timeout = 60
	}
	url := fmt.Sprintf("%s://%s", protocol, in.Host)
	return &RunningHubClient{
		ApiKey:  in.ApiKey,
		Timeout: in.Timeout,
		url:     url,
		httpClient: g.Client().
			SetTimeout(in.Timeout*time.Second).
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

// GetAccountInfo 获取api账户信息
func (c *RunningHubClient) GetAccountInfo(ctx context.Context) (res *GetAccountRes, err error) {
	url := fmt.Sprintf("%s%s", c.url, getAccountInfo)
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
	res = &CreateTaskRes{
		Code: resp.Code,
		Msg:  resp.Msg,
	}
	if resp.Code == 0 {
		if err := json.Unmarshal(resp.Data, &res); err != nil {
			return nil, fmt.Errorf("decode success data fail: %w", err)
		}
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
		FailedReason: &FailedReason{},
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
		FailedReason *FailedReason `json:"failedReason"`
	}
	res.FailedReason.OriginalInfo = string(resp.Data)
	if err := json.Unmarshal(resp.Data, &failData); err != nil {
		return res, nil
	}
	if failData.FailedReason != nil {
		failData.FailedReason.OriginalInfo = string(resp.Data)
		res.FailedReason = failData.FailedReason
	}
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
	switch res.Status {
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

// GetWorkflowJSON 获取工作流JSON
func (c *RunningHubClient) GetWorkflowJSON(ctx context.Context, workflowId string) (res *GetWorkflowJSONRes, err error) {
	if workflowId == "" {
		return nil, gerror.New("workflowId cannot be empty")
	}
	url := fmt.Sprintf("%s%s", c.url, getWorkflowJSON)
	reqBody, _ := json.Marshal(g.Map{
		"workflowId": workflowId,
		"apiKey":     c.ApiKey,
	})
	resp, err := c.doPost(ctx, url, reqBody)
	if err != nil {
		return nil, err
	}
	res = &GetWorkflowJSONRes{
		WorkflowData: make(map[string]WorkflowJSONNodeInfo),
		Code:         resp.Code,
		Msg:          resp.Msg,
	}
	var data struct {
		Prompt string `json:"prompt"`
	}
	if resp.Code == 0 {
		if err := json.Unmarshal(resp.Data, &data); err != nil {
			return nil, fmt.Errorf("decode success prompt string fail: %w", err)
		}
		if err := json.Unmarshal([]byte(data.Prompt), &res.WorkflowData); err != nil {
			return nil, fmt.Errorf("decode success prompt string fail: %w", err)
		}
	}
	return res, nil
}

func (c *RunningHubClient) DownloadWorkflowJsonData(ctx context.Context, in *DownloadWorkflowJSONInput) (filePath string, err error) {
	if in.WorkflowId == "" || in.Client == nil {
		return "", nil
	}
	if in.SaveDir == "" {
		cacheFilePath := filepath.Join("temp", "cacheWorkflowJson")
		os.MkdirAll(cacheFilePath, os.ModePerm)
		in.SaveDir = cacheFilePath
	}
	if in.FileName == "" {
		in.FileName = fmt.Sprintf("workflow_%s.json", in.WorkflowId)
	}
	if !strings.HasSuffix(in.FileName, ".json") {
		in.FileName = fmt.Sprintf("%s.json", in.FileName)
	}
	filePath = filepath.Join(in.SaveDir, in.FileName)
	_, err = os.Stat(filePath)
	if !errors.Is(err, os.ErrNotExist) {
		fmt.Println("file already exists, skipping download")
		return filePath, nil
	}
	result, err := in.Client.GetWorkflowJSON(ctx, in.WorkflowId)
	if err != nil {
		return "", err
	}
	if result.Code != 0 {
		return "", errors.New(result.Msg)
	}
	out, err := json.MarshalIndent(result.WorkflowData, "", "  ")
	if err != nil {
		return "", err
	}
	err = os.WriteFile(filePath, out, os.ModePerm)
	if err != nil {
		return "", err
	}
	return filePath, nil
}

func (c *RunningHubClient) UploadResourceWithURl(ctx context.Context, url string) (res *UploadResourceRes, err error) {
	saveFilePath, err := utility.DownloadFromUrl(url, "")
	if err != nil {
		return nil, fmt.Errorf("download file fail: %w", err)
	}
	defer os.Remove(saveFilePath)
	res, err = c.UploadResource(ctx, saveFilePath)
	if err != nil {
		return nil, fmt.Errorf("upload resource fail: %w", err)
	}
	return res, nil
}

func (c *RunningHubClient) UploadResource(ctx context.Context, filePath string) (res *UploadResourceRes, err error) {
	if !gfile.IsFile(filePath) {
		return nil, errors.New("the filePath does not point to a file")
	}
	url := fmt.Sprintf("%s%s", c.url, uploadResource)
	body := g.Client().
		SetTimeout(c.Timeout*time.Second).
		PostContent(ctx, url, g.Map{
			"apiKey":   c.ApiKey,
			"file":     "@file:" + filePath,
			"fileType": "input",
		})
	var response *RunningHubResponse
	if err = json.Unmarshal([]byte(body), &response); err != nil {
		return nil, err
	}
	if response.Code != 0 {
		return nil, gerror.Newf("UploadResource fail, code: %d, msg: %s", response.Code, response.Msg)
	}
	if err := json.Unmarshal(response.Data, &res); err != nil {
		return nil, fmt.Errorf("decode success data fail: %w", err)
	}
	return res, nil
}

func (c *RunningHubClient) ParseWorkflowPictureInputNode(ctx context.Context, workflowId string) (res []WorkflowNodeInfo, err error) {
	res = []WorkflowNodeInfo{}
	result, err := c.GetWorkflowJSON(ctx, workflowId)
	if err != nil {
		return nil, err
	}
	if result.Code != 0 {
		return nil, errors.New(result.Msg)
	}
	cacheUsedNode := make(map[int]struct{})
	cacheData := make(map[string]*WorkflowNodeInfo)
	for nodeId, v := range result.WorkflowData {
		if ok, nodeInfo := JudgeRunningHubWorkflowNodeIsPictureInputNode(v.ClassType); ok {
			cacheData[nodeId] = &WorkflowNodeInfo{
				NodeId:    nodeId,
				NodeType:  nodeInfo.NodeType,
				FieldName: nodeInfo.FieldName,
			}
		}
		for _, cacheV := range gconv.SliceInt(v.Inputs["image"]) {
			if cacheV > 0 {
				cacheUsedNode[cacheV] = struct{}{}
			}
		}
	}
	res = make([]WorkflowNodeInfo, 0, len(cacheData))
	for nodeId, v := range cacheData {
		nodeIdInt := gconv.Int(nodeId)
		if _, ok := cacheUsedNode[nodeIdInt]; ok {
			res = append(res, *v)
		}
	}
	return res, nil
}
