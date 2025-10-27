package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// 获取文件名（不含后缀）并添加单引号
func getFileNamesWithoutExt(dirPath string) ([]string, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	var namesWithoutExt []string
	for _, entry := range entries {
		if !entry.IsDir() {
			fileName := entry.Name()
			// 获取文件后缀（如 .txt）
			ext := filepath.Ext(fileName)
			// 去除后缀（如果有后缀则截取，无后缀则保持原名称）
			nameWithoutExt := strings.TrimSuffix(fileName, ext)
			// 添加单引号（仅修改了这里的引号类型）
			namesWithoutExt = append(namesWithoutExt, fmt.Sprintf("'%s'", nameWithoutExt))
		}
	}
	return namesWithoutExt, nil
}

func main() {
	targetDir := "../assets/img/touxiang"

	names, err := getFileNamesWithoutExt(targetDir)
	if err != nil {
		fmt.Printf("获取文件名失败：%v\n", err)
		return
	}

	// 用逗号隔开结果
	result := strings.Join(names, ",")
	fmt.Printf("文件夹 %s 下的文件名（不含后缀，带单引号）：%s\n", targetDir, result)
}
