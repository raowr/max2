package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main1() {
	// 指定要处理的目录路径
	root := "../assets" // 请替换为你的实际路径

	// 检查目录是否存在
	if _, err := os.Stat(root); os.IsNotExist(err) {
		fmt.Printf("错误: 目录不存在: %s\n", root)
		return
	}

	fmt.Printf("开始处理目录: %s\n", root)
	fmt.Println("正在扫描文件...")

	var processedCount int

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("访问路径错误 %s: %v\n", path, err)
			return nil // 继续处理其他文件
		}

		// 跳过目录本身，只处理文件
		if info.IsDir() {
			return nil
		}

		// 获取目录和文件名
		dir := filepath.Dir(path)
		filename := info.Name()

		// 检查文件名是否包含空格
		if strings.Contains(filename, " ") {
			// 去掉文件名中的所有空格
			newFilename := strings.ReplaceAll(filename, " #", "")

			// 如果新文件名与原文件名相同，跳过
			if newFilename == filename {
				return nil
			}

			// 构建新的完整路径
			newPath := filepath.Join(dir, newFilename)

			// 检查新路径是否已存在
			if _, err := os.Stat(newPath); err == nil {
				fmt.Printf("警告: 目标文件已存在，跳过: %s\n", newPath)
				return nil
			}

			// 重命名文件
			err := os.Rename(path, newPath)
			if err != nil {
				fmt.Printf("重命名文件失败 %s -> %s: %v\n", path, newPath, err)
				return nil // 继续处理其他文件
			}

			fmt.Printf("重命名: %s -> %s\n", path, newPath)
			processedCount++
		}

		return nil
	})

	if err != nil {
		fmt.Printf("处理过程中发生错误: %v\n", err)
	} else {
		fmt.Printf("处理完成! 共处理了 %d 个文件\n", processedCount)
	}
}
