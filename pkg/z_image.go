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

type ZImageConfig struct {
	Enabled  bool     `json:"enabled,omitempty"`
	BaseURL  string   `json:"base_url"`
	ApiKey   string   `json:"api_key"`
	Model    string   `json:"model"`
	Prompt   string   `json:"prompt"`
	ImageURL []string `json:"image_url,omitempty"`
}

type ZImageResponse struct {
	TaskID string `json:"task_id"`
}

type ZImageTaskResponse struct {
	TaskStatus   string   `json:"task_status"`
	OutputImages []string `json:"output_images"`
}

// ZImageDrawing 造相绘图
func ZImageDrawing(config *ZImageConfig) ([]*string, error) {
	if config.Prompt == "" {
		return nil, fmt.Errorf("绘图提示词为空")
	}
	if config.BaseURL == "" || config.ApiKey == "" {
		return nil, fmt.Errorf("未找到绘图密钥")
	}
	// 设置默认值
	if config.Model == "" {
		config.Model = "Z-Image-Turbo"
	}
	switch config.Model {
	case "Z-Image":
		config.Model = "Tongyi-MAI/Z-Image"
	case "Z-Image-Turbo":
		config.Model = "Tongyi-MAI/Z-Image-Turbo"
	case "Qwen-Image-Edit-2511":
		config.Model = "Qwen/Qwen-Image-Edit-2511"
		// 支持的模型
	default:
		return nil, fmt.Errorf("不支持的造相模型: %s", config.Model)
	}
	// 准备请求体
	requestBody, _ := json.Marshal(map[string]any{
		"model":     config.Model,
		"prompt":    config.Prompt,
		"image_url": config.ImageURL,
	})
	// 创建HTTP请求
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/images/generations", strings.TrimSuffix(config.BaseURL, "/")), bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}
	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.ApiKey)
	req.Header.Set("X-ModelScope-Async-Mode", "true")
	// 发送请求
	client := &http.Client{Timeout: 30 * time.Second}
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
	var zImageResp ZImageResponse
	if err := json.Unmarshal(body, &zImageResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	deadline := time.Now().Add(15 * time.Minute)
	for {
		if time.Now().After(deadline) {
			return nil, fmt.Errorf("造相绘图任务超时")
		}
		taskReq, err := http.NewRequest("GET", fmt.Sprintf("%s/v1/tasks/%s", strings.TrimSuffix(config.BaseURL, "/"), zImageResp.TaskID), nil)
		if err != nil {
			return nil, fmt.Errorf("创建任务查询请求失败: %v", err)
		}
		taskReq.Header.Set("Content-Type", "application/json")
		taskReq.Header.Set("Authorization", "Bearer "+config.ApiKey)
		taskReq.Header.Set("X-ModelScope-Task-Type", "image_generation")

		taskClient := &http.Client{Timeout: 30 * time.Second}
		taskRespRaw, err := taskClient.Do(taskReq)
		if err != nil {
			return nil, fmt.Errorf("任务查询请求失败: %v", err)
		}
		if taskRespRaw.StatusCode < 200 || taskRespRaw.StatusCode >= 300 {
			b, _ := io.ReadAll(taskRespRaw.Body)
			taskRespRaw.Body.Close()
			return nil, fmt.Errorf("任务查询失败: %s", string(b))
		}

		var taskResp ZImageTaskResponse
		if err := json.NewDecoder(taskRespRaw.Body).Decode(&taskResp); err != nil {
			taskRespRaw.Body.Close()
			return nil, fmt.Errorf("解析任务响应失败: %v", err)
		}
		taskRespRaw.Body.Close()

		if taskResp.TaskStatus == "SUCCEED" && len(taskResp.OutputImages) > 0 {
			var urls []*string
			for _, img := range taskResp.OutputImages {
				urls = append(urls, &img)
			}
			return urls, nil
		} else if taskResp.TaskStatus == "FAILED" {
			fmt.Println("Image Generation Failed.")
			break
		}

		time.Sleep(5 * time.Second)
	}

	return nil, fmt.Errorf("造相绘图任务超时")
}
