package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"np/util"
	"os"
)

var packageJSON map[string]string

var tag string
var version string

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
		if version != "" {
			err := util.UpdatePackageVersion(version)
			if err != nil {
				fmt.Printf("Error updating version %s: %s\n", version, err)
				os.Exit(1)
			}
			err = util.RunBuild()
			if err != nil {
				fmt.Printf("Error building version %s: %s\n", version, err)
				os.Exit(1)
			}
		} else {
			if tag != "" {
				currentVersion.IncrementPreRelease(tag)
			} else {

			}
		}

	},
}

func init() {
	rootCmd.AddCommand(publishCmd)
	publishCmd.Flags().StringVarP(&tag, "tag", "t", "", "Tag (alpha, beta, release) of the package to publish")
	publishCmd.Flags().StringVarP(&version, "version", "v", "", "Version of the package to publish")
}
