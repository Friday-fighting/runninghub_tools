package runninghub_tools

import "encoding/json"

type RunningHubClientConfig struct {
	UseHttpReq bool   `json:"use_http_req"`
	Host       string `json:"host"`
	ApiKey     string `json:"api_key"`
}

// RunningHubResponse RunningHub 通用响应结构
type RunningHubResponse struct {
	Code          int             `json:"code"`
	Msg           string          `json:"msg"`
	ErrorMessages interface{}     `json:"errorMessages"`
	Data          json.RawMessage `json:"data"`
}

type CreateTaskReq struct {
	WorkflowId   string      `json:"workflowId"`
	NodeInfoList []*NodeInfo `json:"nodeInfoList"`
	ApiKey       string      `json:"apiKey"`
	WebhookUrl   string      `json:"webhookUrl"`
}

type NodeInfo struct {
	NodeId     string `json:"nodeId"`
	FieldName  string `json:"fieldName"`
	FieldValue string `json:"fieldValue"`
}

type CreateTaskRes struct {
	NetWssUrl  string `json:"netWssUrl"`
	TaskId     string `json:"taskId"`
	TaskStatus string `json:"taskStatus"`
	ClientId   string `json:"clientId"`
	PromptTips string `json:"promptTips"`
}

// GetAccountRes 获取账户信息成功响应数据
type GetAccountRes struct {
	RemainCoins       string `json:"remainCoins"`
	CurrentTaskCounts string `json:"currentTaskCounts"`
	RemainMoney       string `json:"remainMoney"`
	Currency          string `json:"currency"`
	ApiType           string `json:"apiType"`
}

type GetTaskStatusAndResultReq struct {
	TaskId    string
	MaxTries  int
	SleepTime int
}

type GetTaskStatusAndResultRes struct {
	Status       string                                `json:"status"`
	Code         int                                   `json:"code"`
	Msg          string                                `json:"msg"`
	SuccessItems []*SuccessOfGetTaskResultResponseData `json:"success_items"`
	FailedReason *FailedOfGetTaskResultResponseData    `json:"failed_reason"`
}

type GetTaskResultRes struct {
	Code         int                                   `json:"code"`
	Msg          string                                `json:"msg"`
	SuccessItems []*SuccessOfGetTaskResultResponseData `json:"success_items"`
	FailedReason *FailedOfGetTaskResultResponseData    `json:"failed_reason"`
}

// SuccessOfGetTaskResultResponseData 获取任务结果成功响应
type SuccessOfGetTaskResultResponseData struct {
	FileUrl                string  `json:"fileUrl"`
	FileType               string  `json:"fileType"`
	TaskCostTime           string  `json:"taskCostTime"`
	NodeId                 string  `json:"nodeId"`
	ThirdPartyConsumeMoney *string `json:"thirdPartyConsumeMoney"`
	ConsumeMoney           string  `json:"consumeMoney"`
	ConsumeCoins           *string `json:"consumeCoins"`
}

// FailedOfGetTaskResultResponseData 获取任务结果失败响应
type FailedOfGetTaskResultResponseData struct {
	FailedReason struct {
		CurrentOutputs   string `json:"current_outputs"`
		ExceptionType    string `json:"exception_type"`
		CurrentInputs    string `json:"current_inputs"`
		TraceBack        string `json:"traceback"`
		NodeId           string `json:"node_id"`
		ExceptionMessage string `json:"exception_message"`
	} `json:"failedReason"`
}
