package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type JimengRequest struct {
	Model          string    `json:"model"`
	Prompt         string    `json:"prompt"`
	Images         []*string `json:"images,omitempty"`
	Ratio          string    `json:"ratio"`
	Resolution     string    `json:"resolution"`
	Duration       int       `json:"duration,omitempty"`
	FilePaths      []*string `json:"file_paths,omitempty"`
	NegativePrompt string    `json:"negative_prompt,omitempty"`
	SampleStrength float64   `json:"sample_strength,omitempty"`
	ResponseFormat string    `json:"response_format"`
}

type JimengConfig struct {
	BaseURL   string   `json:"base_url"`
	SessionID []string `json:"sessionid"`
	JimengRequest
}

type JimengResponse struct {
	Created int64 `json:"created"`
	Data    []struct {
		URL string `json:"url"`
	} `json:"data"`
}

func intPtr(v int) *int {
	return &v
}

func floatPtr(v float64) *float64 {
	return &v
}

func JimengImageGenerations(config *JimengConfig) ([]*string, error) {
	if config.Prompt == "" {
		return nil, fmt.Errorf("绘图提示词为空")
	}
	if len(config.SessionID) == 0 {
		return nil, fmt.Errorf("未找到绘图密钥")
	}
	// 设置默认值
	if config.Model == "" {
		config.Model = "jimeng-4.0"
	}
	if config.ResponseFormat == "" {
		config.ResponseFormat = "url"
	}
	if config.Ratio == "" {
		config.Ratio = "16:9"
	}
	if config.Resolution == "" {
		config.Resolution = "2k"
	}
	if config.SampleStrength == 0 {
		config.SampleStrength = 0.5
	}
	sessionID := strings.Join(config.SessionID, ",")
	// 准备请求体
	requestBody, err := json.Marshal(config.JimengRequest)
	if err != nil {
		return nil, fmt.Errorf("序列化请求体失败: %v", err)
	}
	// 创建HTTP请求
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/images/generations", config.BaseURL), bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}
	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+sessionID)
	// 发送请求
	client := &http.Client{Timeout: 300 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()
	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}
	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API请求失败，状态码 %d: %s", resp.StatusCode, string(body))
	}
	// 解析响应
	var jimengResp JimengResponse
	if err := json.Unmarshal(body, &jimengResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}
	// 检查是否有生成的图片
	if len(jimengResp.Data) == 0 {
		return nil, fmt.Errorf("未生成任何图片")
	}

	var urls []*string
	for _, data := range jimengResp.Data {
		urls = append(urls, &data.URL)
	}

	return urls, nil
}

func JimengImageCompositions(config *JimengConfig) ([]*string, error) {
	if config.Prompt == "" {
		return nil, fmt.Errorf("绘图提示词为空")
	}
	if len(config.Images) == 0 {
		return nil, fmt.Errorf("输入图像列表不能为空")
	}
	if len(config.SessionID) == 0 {
		return nil, fmt.Errorf("未找到绘图密钥")
	}
	// 设置默认值
	if config.Model == "" {
		config.Model = "jimeng-4.0"
	}
	if config.ResponseFormat == "" {
		config.ResponseFormat = "url"
	}
	if config.Ratio == "" {
		config.Ratio = "16:9"
	}
	if config.Resolution == "" {
		config.Resolution = "2k"
	}
	if config.SampleStrength == 0 {
		config.SampleStrength = 0.5
	}
	sessionID := strings.Join(config.SessionID, ",")
	// 准备请求体
	requestBody, err := json.Marshal(config.JimengRequest)
	if err != nil {
		return nil, fmt.Errorf("序列化请求体失败: %v", err)
	}
	// 创建HTTP请求
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/images/compositions", config.BaseURL), bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}
	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+sessionID)
	// 发送请求
	client := &http.Client{Timeout: 300 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()
	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}
	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API请求失败，状态码 %d: %s", resp.StatusCode, string(body))
	}
	// 解析响应
	var jimengResp JimengResponse
	if err := json.Unmarshal(body, &jimengResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}
	// 检查是否有生成的图片
	if len(jimengResp.Data) == 0 {
		return nil, fmt.Errorf("未生成任何图片")
	}

	var urls []*string
	for _, data := range jimengResp.Data {
		urls = append(urls, &data.URL)
	}

	return urls, nil
}

func JimengVideoGenerations(config *JimengConfig) ([]*string, error) {
	if config.Prompt == "" {
		return nil, fmt.Errorf("绘图提示词为空")
	}
	if len(config.SessionID) == 0 {
		return nil, fmt.Errorf("未找到绘图密钥")
	}
	// 设置默认值
	if config.Model == "" {
		config.Model = "jimeng-video-3.0-fast"
	}
	if config.ResponseFormat == "" {
		config.ResponseFormat = "url"
	}
	if config.Ratio == "" {
		config.Ratio = "16:9"
	}
	if config.Resolution == "" {
		config.Resolution = "720p"
	}
	if config.Duration == 0 {
		config.Duration = 5
	}
	sessionID := strings.Join(config.SessionID, ",")
	// 准备请求体
	requestBody, err := json.Marshal(config.JimengRequest)
	if err != nil {
		return nil, fmt.Errorf("序列化请求体失败: %v", err)
	}
	// 创建HTTP请求
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/videos/generations", config.BaseURL), bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}
	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+sessionID)
	// 发送请求
	client := &http.Client{Timeout: 300 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()
	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}
	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API请求失败，状态码 %d: %s", resp.StatusCode, string(body))
	}
	// 解析响应
	var jimengResp JimengResponse
	if err := json.Unmarshal(body, &jimengResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}
	// 检查是否有生成的视频
	if len(jimengResp.Data) == 0 {
		return nil, fmt.Errorf("未生成任何视频")
	}

	var urls []*string
	for _, data := range jimengResp.Data {
		urls = append(urls, &data.URL)
	}

	return urls, nil
}
