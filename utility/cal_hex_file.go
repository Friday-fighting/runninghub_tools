package utility

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
)

func GetLocalFileMd5Hex(filepath string) (md5Hex string, err error) {
	file, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// 创建 MD5 hash
	hash := md5.New()

	// 流式拷贝到 hash
	if _, err := io.Copy(hash, file); err != nil {
		panic(err)
	}

	// 计算结果
	md5Bytes := hash.Sum(nil)

	// 转为 hex 字符串
	md5Hex = hex.EncodeToString(md5Bytes)
	return md5Hex, nil
}
