package runninghub_client_utils

type NodeMeta struct {
	ClassType string `json:"class_type"`
	FieldName string `json:"field_name"`
	NodeType  string `json:"node_type"`
}

var RunningHubWorkflowPictureInputNodeInfo = map[string]NodeMeta{
	// 上传图片（base64）节点
	"LoadImageFromBase64": {
		ClassType: "LoadImageFromBase64",
		FieldName: "data",
		NodeType:  "base64",
	},
	"easy loadImageBase64": {
		ClassType: "easy loadImageBase64",
		FieldName: "base64_data",
		NodeType:  "base64",
	},
	// 上传图片（文件）节点
	"LoadImageHDR": {
		ClassType: "LoadImageHDR",
		FieldName: "image",
		NodeType:  "file",
	},
	"LoadImageMask": {
		ClassType: "LoadImageMask",
		FieldName: "image",
		NodeType:  "file",
	},
	"LoadHDRImage": {
		ClassType: "LoadHDRImage",
		FieldName: "image",
		NodeType:  "file",
	},
	"LoadPILImage": {
		ClassType: "LoadPILImage",
		FieldName: "image",
		NodeType:  "file",
	},
	"LoadImageOutput": {
		ClassType: "LoadImageOutput",
		FieldName: "image",
		NodeType:  "file",
	},
	"LoadImage": {
		ClassType: "LoadImage",
		FieldName: "image",
		NodeType:  "file",
	},
	"LoadImageReturnFilename": {
		ClassType: "LoadImageReturnFilename",
		FieldName: "image",
		NodeType:  "file",
	},
	"LoadImageMW": {
		ClassType: "LoadImageMW",
		FieldName: "image",
		NodeType:  "file",
	},
	"ZML_LoadImage": {
		ClassType: "ZML_LoadImage",
		FieldName: "图像",
		NodeType:  "file",
	},
	"MuyeLoadImage": {
		ClassType: "MuyeLoadImage",
		FieldName: "image",
		NodeType:  "file",
	},
	"Load image with metadata [Crystools]": {
		ClassType: "Load image with metadata [Crystools]",
		FieldName: "image",
		NodeType:  "file",
	},
	"LoadImage //Inspire": {
		ClassType: "LoadImage //Inspire",
		FieldName: "image",
		NodeType:  "file",
	},
	// 从 URL 加载图像节点
	"LoadImageFromUrl": {
		ClassType: "LoadImageFromUrl",
		FieldName: "image",
		NodeType:  "url",
	},
	"LoadImageAsMaskFromUrl": {
		ClassType: "LoadImageAsMaskFromUrl",
		FieldName: "image",
		NodeType:  "url",
	},
	"Load Image From Url (mtb)": {
		ClassType: "Load Image From Url (mtb)",
		FieldName: "url",
		NodeType:  "url",
	},
	"LoadImageFromURL": {
		ClassType: "LoadImageFromURL",
		FieldName: "url",
		NodeType:  "url",
	},
	"LoadImagesFromURL": {
		ClassType: "LoadImagesFromURL",
		FieldName: "url",
		NodeType:  "url",
	},
	"Light-Tool: LoadImageFromURL": {
		ClassType: "Light-Tool: LoadImageFromURL",
		FieldName: "url",
		NodeType:  "url",
	},
}

func JudgeRunningHubWorkflowNodeIsPictureInputNode(classType string) (sign bool, nodeInfo *NodeMeta) {
	if nodeInfo, ok := RunningHubWorkflowPictureInputNodeInfo[classType]; ok {
		return true, &nodeInfo
	}
	return false, nil
}
