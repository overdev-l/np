package util

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type PackageJSON struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type Version struct {
	Major            int    // 主要版本
	Minor            int    // 次要版本
	Patch            int    // 修订版本
	PreRelease       string // 预发布版本
	PreReleaseNumber int    // 预发布版本号
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
	// 转换匹配到的字段
	major, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil, fmt.Errorf("error converting major version: %w", err)
	}

	minor, err := strconv.Atoi(matches[2])
	if err != nil {
		return nil, fmt.Errorf("error converting minor version: %w", err)
	}

	patch, err := strconv.Atoi(matches[3])
	if err != nil {
		return nil, fmt.Errorf("error converting patch version: %w", err)
	}

	preReleaseNumber := 0
	if matches[5] != "" {
		preReleaseNumber, err = strconv.Atoi(matches[5])
		if err != nil {
			return nil, fmt.Errorf("error converting pre-release number: %w", err)
		}
	}
	return &Version{
		Major:            major,
		Minor:            minor,
		Patch:            patch,
		PreRelease:       matches[4],
		PreReleaseNumber: preReleaseNumber,
	}, nil
}

func (v *Version) IncrementPreRelease(tag string) {
	if v.PreRelease == tag {
		v.PreReleaseNumber++
	} else {
		v.PreRelease = tag
		v.PreReleaseNumber = 1
	}
}

func (v *Version) UpdatePackageReleaseVersion(preReleaseVersion int) {
	v.PreReleaseNumber = preReleaseVersion
}

func (v Version) String() string {
	version := fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
	if v.PreRelease != "" {
		version += "-" + v.PreRelease
		if v.PreReleaseNumber > 0 {
			version += "." + strconv.Itoa(v.PreReleaseNumber)
		}
	}
	return version
}

func RunBuild() error {
	cmd := exec.Command("npm", "run", "build")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return err
	}

	fmt.Printf("Output:\n%s\n", output)
	return nil
}

func RunCommand(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
