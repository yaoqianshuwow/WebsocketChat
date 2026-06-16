package node

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/bytedance/gopkg/util/logger"

	"ppt-agent/biz/graph/state"
	"ppt-agent/pkg/myfile"
)

// GenerateImage 调用 agnes-image-2.0-flash 生成单张图片，返回本地路径
func GenerateImage(ctx context.Context, description, pageTitle, subDir string) (*state.ImageInfo, error) {
	if contentImageModel == nil {
		return nil, fmt.Errorf("image model not initialized")
	}

	payload := map[string]any{
		"model":  contentImageModel.ModelName,
		"prompt": fmt.Sprintf("Professional presentation image: %s", description),
		"n":      1,
		"size":   "1024x1024",
		"response_format": "b64_json",
	}
	body, _ := json.Marshal(payload)

	url := fmt.Sprintf("%s/images/generations", contentImageModel.BaseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+contentImageModel.APIKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API请求失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API错误(%d): %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		Data []struct {
			B64JSON string `json:"b64_json"`
			URL     string `json:"url"`
		} `json:"data"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}
	if len(result.Data) == 0 {
		return nil, fmt.Errorf("API未返回图片")
	}

	// 保存图片
	root, _ := myfile.GetProjectRoot()
	imgDir := filepath.Join(root, "output", "images", subDir, fmt.Sprintf("%d", time.Now().UnixMilli()))
	os.MkdirAll(imgDir, 0755)

	fileName := fmt.Sprintf("%s.png", pageTitle)
	if len(fileName) > 80 {
		fileName = fileName[:80] + ".png"
	}
	filePath := filepath.Join(imgDir, fileName)

	var imgData []byte
	if result.Data[0].B64JSON != "" {
		imgData, err = base64.StdEncoding.DecodeString(result.Data[0].B64JSON)
		if err != nil {
			return nil, fmt.Errorf("base64解码失败: %w", err)
		}
	} else if result.Data[0].URL != "" {
		imgResp, err := http.Get(result.Data[0].URL)
		if err != nil {
			return nil, fmt.Errorf("下载图片失败: %w", err)
		}
		defer imgResp.Body.Close()
		imgData, err = io.ReadAll(imgResp.Body)
		if err != nil {
			return nil, fmt.Errorf("读取图片失败: %w", err)
		}
	} else {
		return nil, fmt.Errorf("API未返回图片数据")
	}

	if err := os.WriteFile(filePath, imgData, 0644); err != nil {
		return nil, fmt.Errorf("保存图片失败: %w", err)
	}

	logger.Infof("图片生成成功: %s (%d bytes)", filePath, len(imgData))
	return &state.ImageInfo{
		PageTitle:   pageTitle,
		Description: description,
		LocalPath:   filePath,
	}, nil
}

// searchFallbackImage 搜索兜底：当 AI 生成失败时，返回描述文字作为占位
func searchFallbackImage(ctx context.Context, query, pageTitle string) (*state.ImageInfo, error) {
	root, _ := myfile.GetProjectRoot()
	fallbackDir := filepath.Join(root, "output", "images", "fallback")
	os.MkdirAll(fallbackDir, 0755)

	// 将搜索词写入文本占位文件
	placeholder := fmt.Sprintf("search_fallback_%s.txt", pageTitle)
	if len(placeholder) > 80 {
		placeholder = placeholder[:80] + ".txt"
	}
	filePath := filepath.Join(fallbackDir, placeholder)
	content := fmt.Sprintf("Search query: %s\nPage: %s\n(Image generation failed, please replace manually)", query, pageTitle)
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return nil, fmt.Errorf("写入兜底文件失败: %w", err)
	}

	logger.Infof("搜索兜底已创建: %s", filePath)
	return &state.ImageInfo{
		PageTitle:   pageTitle,
		Description: query,
		LocalPath:   filePath,
	}, nil
}
