package runninghub_tools

import "strings"

type ErrorInfo struct {
	Code         int           `json:"code"`
	Sign         string        `json:"sign"`
	Msg          string        `json:"msg"`
	SubCode      int           `json:"sub_code"`
	FailedReason *FailedReason `json:"failed_reason"`
}

type errorInfo struct {
	Code int    `json:"code"`
	Sign string `json:"sign"`
	Msg  string `json:"msg"`
}

var errorInfoMap = map[int]errorInfo{
	301: {
		Code: 301,
		Sign: "PARAMS_INVALID",
		Msg:  "请求中包含非法或缺失的参数",
	},
	380: {
		Code: 380,
		Sign: "WORKFLOW_NOT_EXISTS",
		Msg:  "指定的工作流不存在",
	},
	412: {
		Code: 412,
		Sign: "TOKEN_INVALID",
		Msg:  "API接口路径拼写错误",
	},
	415: {
		Code: 415,
		Sign: "TASK_INSTANCE_MAXED",
		Msg:  "独占型 API 当前可用的实例/机器数不足",
	},
	416: {
		Code: 416,
		Sign: "TASK_CREATE_FAILED_BY_NOT_ENOUGH_WALLET",
		Msg:  "钱包余额不足",
	},
	421: {
		Code: 421,
		Sign: "TASK_QUEUE_MAXED",
		Msg:  "共享型 API 的并发数已达到用户上限",
	},
	423: {
		Code: 423,
		Sign: "TASK_NOT_FOUNED",
		Msg:  "未找到指定任务",
	},
	433: {
		Code: 433,
		Sign: "VALIDATE_PROMPT_FAILED",
		Msg:  "工作流合法性校验未通过（包含 prompt 与节点配置校验）",
	},
	435: {
		Code: 435,
		Sign: "TASK_USER_EXCLAPI_INSTANCE_NOT_FOUND",
		Msg:  "未找到任务用户API实例",
	},
	436: {
		Code: 436,
		Sign: "TASK_USER_EXCLAPI_REQUIRED",
		Msg:  "独占会员到期",
	},
	801: {
		Code: 801,
		Sign: "APIKEY_UNSUPPORTED_FREE_USER",
		Msg:  "免费用户不支持 API Key",
	},
	802: {
		Code: 802,
		Sign: "APIKEY_UNAUTHORIZED",
		Msg:  "API Key 未授权/已失效",
	},
	803: {
		Code: 803,
		Sign: "APIKEY_INVALID_NODE_INFO",
		Msg:  "传入的 nodeInfoList 与绑定的工作流不匹配",
	},
	804: {
		Code: 804,
		Sign: "APIKEY_TASK_IS_RUNNING",
		Msg:  "任务正在运行中",
	},
	805: {
		Code: 805,
		Sign: "APIKEY_TASK_STATUS_ERROR",
		Msg:  "任务状态异常",
	},
	806: {
		Code: 806,
		Sign: "APIKEY_USER_NOT_FOUND",
		Msg:  "未找到对应用户",
	},
	807: {
		Code: 807,
		Sign: "APIKEY_TASK_NOT_FOUND",
		Msg:  "未找到对应任务",
	},
	808: {
		Code: 808,
		Sign: "APIKEY_UPLOAD_FAILED",
		Msg:  "文件上传失败",
	},
	809: {
		Code: 809,
		Sign: "APIKEY_FILE_SIZE_EXCEEDED",
		Msg:  "文件大小超出限制",
	},
	810: {
		Code: 810,
		Sign: "WORKFLOW_NOT_SAVED_OR_NOT_RUNNING",
		Msg:  "用户未保存工作流或未在平台运行直接调用api",
	},
	811: {
		Code: 811,
		Sign: "CORPAPIKEY_INVALID",
		Msg:  "企业版 API Key 无效",
	},
	812: {
		Code: 812,
		Sign: "CORPAPIKEY_INSUFFICIENT_FUNDS",
		Msg:  "企业版余额不足",
	},
	813: {
		Code: 813,
		Sign: "APIKEY_TASK_IS_QUEUED",
		Msg:  "任务已排队等待执行",
	},
	901: {
		Code: 901,
		Sign: "WEBAPP_NOT_EXISTS",
		Msg:  "WebApp 不存在",
	},
	500: {
		Code: 500,
		Sign: "UNKNOWN_ERROR",
		Msg:  "未知错误(未被显式捕获的异常)",
	},
}

func GetErrorInfo(code int, failReason *FailedReason) (res *ErrorInfo) {
	res = &ErrorInfo{
		Code:         code,
		SubCode:      code,
		FailedReason: failReason,
	}
	if code == 805 {
		res.Msg = "任务状态异常"
		res.Sign = "APIKEY_TASK_STATUS_ERROR"
		if strings.Contains(res.FailedReason.ExceptionMessage, "Porn") {
			res.SubCode = 805001
			return res
		}
		if strings.Contains(res.FailedReason.ExceptionMessage, "显存告警") {
			res.SubCode = 805002
			return res
		}
		res.SubCode = 805000
		return res
	}
	if value, ok := errorInfoMap[code]; ok {
		res.Sign = value.Sign
		res.Msg = value.Msg
		return res
	}
	res.Msg = failReason.TraceBack
	res.Sign = failReason.ExceptionMessage
	return res
}
