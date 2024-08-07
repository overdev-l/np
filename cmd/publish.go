package cmd

import (
	"bytes"
	"fmt"
	"github.com/spf13/cobra"
	"np/util"
	"os"
	"os/exec"
)

var packageJSON map[string]string

var tag string
var version string
var preReleaseNumber int

var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "Publish an package to a registry",
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		pkg, err := util.GetPackageJSON()
		if err != nil {
			os.Exit(1)
		}
		packageJSON = pkg
	},
	Run: func(cmd *cobra.Command, args []string) {
		var currentVersion *util.Version
		var stringVersion string
		if version == "" {
			stringVersion = version
		} else {
			stringVersion = packageJSON["version"]
		}
		parseVersion, err := util.ParseVersion(stringVersion)
		if err != nil {
			fmt.Printf("Error parsing version %s: %s\n", version, err)
			os.Exit(1)
		}
		currentVersion = parseVersion
		// 如果显式指定了版本号(--version)则 --tag 无效
		// 指定了 --tag 才会自动增加 PreRelease的版本号
		// 无 --tag 自动新增修订版本号
		// 不需要自动新增则需显式 --version 指定版本号
		if tag != "" && version == "" {
			if preReleaseNumber < 0 {
				fmt.Printf("Error: pre-release number not specified\n")
				os.Exit(1)
			}
			currentVersion.IncrementPreRelease(tag)
			currentVersion.UpdatePackageReleaseVersion(preReleaseNumber)
		} else {
			currentVersion.Patch++
		}
		err = util.UpdatePackageVersion(currentVersion.String())
		if err != nil {
			os.Exit(1)
		}
		err = util.RunBuild()
		if err != nil {
			fmt.Printf("Error building version %s: %s\n", version, err)
			os.Exit(1)
		}
		err = util.RunCommand("git", "add", ".")
		if err != nil {
			fmt.Printf("Error adding version %s: %s\n", version, err)
			os.Exit(1)
		}
		err = util.RunCommand("git", "commit", "-m", fmt.Sprintf("Added version %s", version))
		if err != nil {
			fmt.Printf("Error adding version %s: %s\n", version, err)
			os.Exit(1)
		}
		err = util.RunCommand("git", "push", "origin", "HEAD")
		if err != nil {
			fmt.Printf("Error adding version %s: %s\n", version, err)
			os.Exit(1)
		}
		config, err := util.GetConfig()
		if err != nil {
			fmt.Printf("Error getting config: %s\n", err)
			os.Exit(1)
		}
		if config["registry"] != "" {
			cmd := exec.Command("npm", "config", "set", "registry", config["registry"])
			_, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Printf("Error adding version %s: %s\n", version, err)
				os.Exit(1)
			}
		}
		loginCommand := exec.Command("npm", "login")
		var stdin bytes.Buffer
		stdin.WriteString(fmt.Sprintf("%s\n%s\n%s\n", config["username"], config["password"], ""))
		loginCommand.Stdin = &stdin
		loginOutput, err := loginCommand.CombinedOutput()
		if err != nil {
			fmt.Printf("Error adding version %s: %s\n", version, err)
			os.Exit(1)
		}
		fmt.Printf("%s\n", string(loginOutput))

		if currentVersion.PreRelease != "" {
			fmt.Printf("Pre-release version: %s\n", currentVersion.PreRelease)
			publishCommand := exec.Command("npm", "publish", "--tag", currentVersion.PreRelease)
			publishOutput, err := publishCommand.CombinedOutput()
			if err != nil {
				fmt.Printf("Error adding version %s: %s\n", version, err)
				os.Exit(1)
			}
			fmt.Printf("%s\n", string(publishOutput))
		} else {
			fmt.Printf("No pre-release version specified.\n")
			publishCommand := exec.Command("npm", "publish")
			publishOutput, err := publishCommand.CombinedOutput()
			if err != nil {
				fmt.Printf("Error adding version %s: %s\n", version, err)
				os.Exit(1)
			}
			fmt.Printf("%s\n", string(publishOutput))
		}
	},
}

func init() {
	rootCmd.AddCommand(publishCmd)
	publishCmd.Flags().StringVarP(&tag, "tag", "t", "", "Tag (alpha, beta, release) of the package to publish")
	publishCmd.Flags().IntVarP(&preReleaseNumber, "preReleaseNumber", "p", -1, "Number of pre-release versions to publish")
	publishCmd.Flags().StringVarP(&version, "version", "v", "", "Version of the package to publish")
}
