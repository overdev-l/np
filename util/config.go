package util

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type PackageJSON struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type Version struct {
	Major            string // 主要版本
	Minor            string // 次要版本
	Patch            string // 修订版本
	PreRelease       string // 预发布版本
	PreReleaseNumber string // 预发布版本号
}

func NpConfig() error {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		return err
	}
	npconfigPath := filepath.Join(home, ".npconfig")
	if _, err := os.Stat(npconfigPath); err == nil {
		return nil
	} else if os.IsNotExist(err) {
		file, err := os.Create(npconfigPath)
		if err != nil {
			fmt.Println("Error creating .npconfig file:", err)
			return err
		}
		defer file.Close()
	} else {
		return err
	}
	return nil
}

func GetConfig() (map[string]string, error) {
	var result = make(map[string]string)
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		return result, err
	}
	npconfigPath := filepath.Join(home, ".npconfig")

	file, err := os.Open(npconfigPath)
	if err != nil {
		fmt.Println(err)
		return result, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			result[key] = value
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		return result, err
	}
	return result, nil
}

func WriteConfig(key, value string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		return err
	}
	filePath := filepath.Join(home, ".npconfig")
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer file.Close()
	var lines []string
	exists := false
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			currentKey := strings.TrimSpace(parts[0])
			currentValue := strings.TrimSpace(parts[1])
			if currentKey == key {
				lines = append(lines, fmt.Sprintf("%s=%s", key, value))
				exists = true
			} else {
				lines = append(lines, fmt.Sprintf("%s=%s", key, currentValue))
			}
		}
	}
	if !exists {
		lines = append(lines, fmt.Sprintf("%s=%s", key, value))
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	if err := file.Truncate(0); err != nil {
		return err
	}
	if _, err := file.WriteString(strings.Join(lines, "\n")); err != nil {
		return err
	}
	if err := file.Sync(); err != nil {
		return err
	}
	return nil
}

func GetPackageJSON() (map[string]string, error) {
	fileContent, err := os.ReadFile("package.json")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	var pkg PackageJSON
	err = json.Unmarshal(fileContent, &pkg)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	result := make(map[string]string)

	return result, nil
}

func UpdatePackageVersion(value string) error {
	fileContent, err := os.ReadFile("package.json")
	if err != nil {
		fmt.Println(err)
		return err
	}
	var pkg PackageJSON
	err = json.Unmarshal(fileContent, &pkg)
	if err != nil {
		fmt.Println(err)
		return err
	}
	pkg.Version = value
	updatedContent, err := json.MarshalIndent(pkg, "", "  ")
	if err != nil {
		fmt.Println(err)
		return err
	}
	if err := os.WriteFile("package.json", updatedContent, 0644); err != nil {
		return fmt.Errorf("error writing file: %w", err)
	}
	return nil
}

func ParseVersion(version string) (*Version, error) {
	re := regexp.MustCompile(`^(\d+)\.(\d+)\.(\d+)(?:-(alpha|beta)\.(\d+))?$`)
	matches := re.FindStringSubmatch(version)

	if matches == nil {
		return nil, fmt.Errorf("invalid version format")
	}

	return &Version{
		Major:            matches[1],
		Minor:            matches[2],
		Patch:            matches[3],
		PreRelease:       matches[4],
		PreReleaseNumber: matches[5],
	}, nil
}
