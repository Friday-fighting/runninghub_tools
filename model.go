package runninghub_tools

type NodeInfo struct {
	NodeId     string `json:"nodeId"`
	FieldName  string `json:"fieldName"`
	FieldValue string `json:"fieldValue"`
}

type CreateTaskRequestInfo struct {
	WorkflowId   string      `json:"workflowId"`
	NodeInfoList []*NodeInfo `json:"nodeInfoList"`
	ApiKey       string      `json:"apiKey"`
	WebhookUrl   string      `json:"webhookUrl"`
}

// 获取账户信息响应
type RunningHubAccountResponse struct {
	ApiType           string  `json:"api_type"`
	Currency          string  `json:"currency"`
	CurrentTaskCounts int     `json:"current_task_counts"`
	RemainCoins       float64 `json:"remain_coins"`
	RemainMoney       float64 `json:"remain_money"`
}

// RunningHub 通用响应结构
type RunningHubResponse[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}

func (r RunningHubResponse[T]) OK() bool { return r.Code == 0 }

// 获取账户信息成功响应数据
type RunningHubAccountSuccessResponseData struct {
	RemainCoins       string `json:"remainCoins"`
	CurrentTaskCounts string `json:"currentTaskCounts"`
}

// 创建任务成功响应
type RunningHubCreateTaskSuccessResponse struct {
	Code int                                      `json:"code"`
	Msg  string                                   `json:"msg"`
	Data *RunningHubCreateTaskSuccessResponseData `json:"data"`
}

type RunningHubCreateTaskSuccessResponseData struct {
	NetWssUrl  string `json:"netWssUrl"`
	TaskId     string `json:"taskId"`
	TaskStatus string `json:"taskStatus"`
	ClientId   string `json:"clientId"`
	PromptTips string `json:"promptTips"`
}

// 获取任务状态成功响应
type RunningHubGetTaskStatusSuccessResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data string `json:"data"`
}

// 获取任务结果成功响应

type RunningHubGetTaskResultSuccessResponseData struct {
	PreviewUrl   string `json:"previewUrl"`
	FileUrl      string `json:"fileUrl"`
	FileType     string `json:"fileType"`
	TaskCostTime string `json:"taskCostTime"`
	TaskStatus   string `json:"taskStatus"`
}

type RunningHubGetTaskResultFailedResponseData struct {
	FailedReason map[string]map[string]string `json:"failedReason"`
}
