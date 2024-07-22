package util

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// 将配置文件转换为map
func ReadConfig() (map[string]string, error) {
	kv := make(map[string]string)
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		return kv, err
	}
	filename := ".npconfig"
	filepath := filepath.Join(home, filename)

	file, err := os.Open(filepath)
	if err != nil {
		if os.IsNotExist(err) {
			// 文件不存在，创建它
			file, err = os.Create(filepath)
			if err != nil {
				fmt.Println("Error creating file:", err)
				return kv, err
			}
			defer file.Close()
			// 返回空map
			return kv, nil
		}
		fmt.Println("Error opening file:", err)
		return kv, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			kv[key] = value
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return kv, err
	}
	return kv, nil
}

func WriteConfig(config map[string]string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		return err
	}
	filename := ".npconfig"
	filepath := filepath.Join(home, filename)
	_, err = os.Stat(filepath)
	if os.IsNotExist(err) {
		// 文件不存在，创建它
		file, err := os.Create(filepath)
		if err != nil {
			fmt.Println("Error creating file:", err)
			return nil
		}
		defer file.Close()
	} else if err != nil {
		fmt.Println("Error checking file:", err)
		return nil
	}
	outputFile, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return err
	}
	defer outputFile.Close()

	writer := bufio.NewWriter(outputFile)
	defer writer.Flush()

	for key, value := range config {
		line := fmt.Sprintf("%s=%s\n", key, value)
		_, err := writer.WriteString(line)
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return err
		}
	}

	return nil
}
