package cmd

import (
	"github.com/spf13/cobra"
	"np/util"
	"os"
)

var packageJSON map[string]string

var tag string
var auto bool

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
}

func init() {
	rootCmd.AddCommand(publishCmd)
	publishCmd.Flags().StringVarP(&tag, "tag", "t", "", "Tag (alpha, beta, release) of the package to publish")
	publishCmd.Flags().BoolVarP(&auto, "auto", "a", false, "Automatically generate an auto-generated package")
}
