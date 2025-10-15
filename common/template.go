package common

import (
	"io"
	"os"
	"path/filepath"
	"strings"
)

const templatesDir = "templates"

// LoadTemplate 只读取运行目录下的 templates 目录，防止其他文件内容泄漏
func LoadTemplate(templateName string) ([]byte, error) {
	// 清理路径，防止目录遍历攻击
	cleanTemplateName := filepath.Clean(templateName)

	// 检查是否尝试访问父目录
	if strings.HasPrefix(cleanTemplateName, "..") || strings.Contains(cleanTemplateName, string(filepath.Separator)+".."+string(filepath.Separator)) {
		return nil, NewFileNotFoundError(templateName) // 拒绝包含父目录的路径
	}

	// 构建完整路径，确保只从 templates 目录读取
	fullPath := filepath.Join(templatesDir, cleanTemplateName)

	if _, err := os.Stat(fullPath); err == nil {
		file, err := os.Open(fullPath)
		if err != nil {
			return nil, err
		}
		defer func(file *os.File) {
			if file != nil {
				_ = file.Close()
			}
		}(file)
		result, err := io.ReadAll(file)
		if err != nil {
			return nil, err
		}
		return result, nil
	}
	return nil, NewFileNotFoundError(templateName)
}
