package helper

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/cezan1/GaiO/internal/config"
)

// 修改类型定义，以实现按 requestid 为 key，下面为 map[时间]value 结构
type jsonData map[string]map[string]string

// 获取当天日期命名的文件路径
func getTodayFilePath() (string, error) {
	today := time.Now().Format("20060102")
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	jsonDir := filepath.Join(currentDir, "jsonFile")
	if err := os.MkdirAll(jsonDir, 0755); err != nil {
		return "", err
	}
	return filepath.Join(jsonDir, fmt.Sprintf("%s.txt", today)), nil
}

// WriteRequestData 以 requestid 为 key 写入 JSON 形式数据，按时间记录 value
func WriteRequestData(requestID string, data string) error {
	filePath, err := getTodayFilePath()
	if err != nil {
		return err
	}

	// 读取现有数据
	existingData, err := readFileData(filePath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	// 初始化当前 requestID 的数据存储结构
	if existingData[requestID] == nil {
		existingData[requestID] = make(map[string]string)
	}

	// 使用当前时间作为键
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	existingData[requestID][currentTime] = data

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(existingData)
}

// readFileData 读取文件中的数据
func readFileData(filePath string) (jsonData, error) {
	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return make(jsonData), nil
		}
		return nil, err
	}
	defer file.Close()

	var data jsonData
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		if err.Error() == "EOF" {
			return make(jsonData), nil
		}
		return nil, err
	}
	return data, nil
}

// GetDataByRequestID 读取以 requestid 为 key 的数据
func GetDataByRequestID(requestID string) (map[string]string, error) {
	filePath, err := getTodayFilePath()
	if err != nil {
		return nil, err
	}
	data, err := readFileData(filePath)
	if err != nil {
		return nil, err
	}

	return data[requestID], nil
}

// 根间排序并去重据时
func GetDataByRequestIDDesc(requestID string) ([]string, error) {
	filePath, err := getTodayFilePath()
	if err != nil {
		return nil, err
	}

	data, err := readFileData(filePath)
	if err != nil {
		return nil, err
	}

	requestData, exists := data[requestID]
	if !exists {
		return []string{}, nil
	}

	// 提取所有时间戳
	timestamps := make([]string, 0, len(requestData))
	for ts := range requestData {
		timestamps = append(timestamps, ts)
	}

	// 将时间戳转换为 time.Time 类型并排序
	timeSlice := make([]time.Time, len(timestamps))
	for i, ts := range timestamps {
		t, err := time.Parse("2006-01-02 15:04:05", ts)
		if err != nil {
			return nil, err
		}
		timeSlice[i] = t
	}
	sort.Slice(timeSlice, func(i, j int) bool {
		return timeSlice[i].After(timeSlice[j])
	})

	// 按降序提取数据并去重
	result := make([]string, 0)
	seen := make(map[string]bool)
	for _, t := range timeSlice {
		ts := t.Format("2006-01-02 15:04:05")
		value := requestData[ts]
		if !seen[value] {
			seen[value] = true
			result = append(result, value)
		}
	}
	if len(result) > config.MaxQuestionNum {
		result = result[:config.MaxQuestionNum]
	}
	return result, nil
}
