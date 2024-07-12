package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "np",
	Short: "np 是一个npm package 发布工具",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func Execute() {
	rootCmd.Execute()
}
