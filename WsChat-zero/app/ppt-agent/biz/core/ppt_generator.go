package core

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
	"unicode"

	"ppt-agent/pkg/myfile"
)

type PptGenerator struct {
	outputDir string
}

func NewPptGenerator() *PptGenerator {
	root, err := myfile.GetProjectRoot()
	if err != nil {
		panic(fmt.Errorf("获取项目根目录失败: %w", err))
	}
	outputDir := filepath.Join(root, "output")
	if err := myfile.EnsureDir(outputDir); err != nil {
		panic(fmt.Errorf("创建输出目录失败: %w", err))
	}
	return &PptGenerator{outputDir: outputDir}
}

// sanitizeFilename 清理文件名，保留中文和常见字符，移除非法字符
func sanitizeFilename(name string) string {
	var b strings.Builder
	for _, r := range name {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.Is(unicode.Scripts["Han"], r) {
			b.WriteRune(r)
		} else if r == ' ' || r == '-' || r == '_' {
			b.WriteRune(r)
		} else if r == '：' || r == '·' || r == '、' || r == '（' || r == '）' {
			b.WriteRune(r)
		}
	}
	result := strings.TrimSpace(b.String())
	if result == "" {
		result = "PPT"
	}
	return result
}

func (g *PptGenerator) CreatePptFile(ctx context.Context, slidesJSON string, style string, topic string) (string, error) {
	tmpFile := filepath.Join(g.outputDir, "slides_temp.json")
	if err := os.WriteFile(tmpFile, []byte(slidesJSON), 0644); err != nil {
		return "", fmt.Errorf("写入临时 slides 文件失败: %w", err)
	}
	defer os.Remove(tmpFile)

	now := time.Now()
	dateStr := now.Format("2006-01-02_150405")
	namePart := sanitizeFilename(topic)
	runes := []rune(namePart)
	if len(runes) > 40 {
		namePart = string(runes[:40])
	}
	fileName := fmt.Sprintf("%s_%s.pptx", namePart, dateStr)
	outputPath := filepath.Join(g.outputDir, fileName)

	root, err := myfile.GetProjectRoot()
	if err != nil {
		return "", fmt.Errorf("获取项目根目录失败: %w", err)
	}

	scriptPath := filepath.Join(root, "scripts", "generate_ppt.py")
	cmd := exec.CommandContext(ctx, "py", "-3",
		scriptPath,
		"--input", tmpFile,
		"--output", outputPath,
		"--style", style,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("生成 PPT 失败: %s: %w", string(output), err)
	}

	return outputPath, nil
}

func (g *PptGenerator) GetOutputDir() string {
	return g.outputDir
}
