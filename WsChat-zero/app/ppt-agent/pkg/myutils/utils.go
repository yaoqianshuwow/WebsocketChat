package myutils

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

// CloseBrowser 关闭浏览器（占位函数，实际项目中可扩展）
func CloseBrowser() {
	fmt.Println("正在清理资源...")
}

// OpenBrowser 打开浏览器
func OpenBrowser(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	default:
		return fmt.Errorf("不支持的操作系统: %s", runtime.GOOS)
	}
	return cmd.Start()
}

// GetEnvOrDefault 获取环境变量，不存在则返回默认值
func GetEnvOrDefault(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}
