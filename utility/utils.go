package utility

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	runninghubClientUtil "github.com/Friday-fighting/runninghub_tools/runninghub_client_utils"
	"os"
	"path/filepath"
	"strings"
)

type DownloadWorkflowJSONInput struct {
	WorkflowId string
	Client     *runninghubClientUtil.RunningHubClient
	SaveDir    string
	FileName   string
}

func DownloadWorkflowJsonData(ctx context.Context, in *DownloadWorkflowJSONInput) (filePath string, err error) {
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
