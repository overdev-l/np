package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"np/util"
	"os"
)

var name string
var password string
var registry string

var configCmd = &cobra.Command{
	Use:   "set",
	Short: "Set configuration",
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		err := util.NpConfig()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		// 如果名称不为空，则写入配置文件
		if name != "" {
			err := util.WriteConfig("username", name)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
		if password != "" {
			err := util.WriteConfig("password", password)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
		if registry != "" {
			err := util.WriteConfig("registry", registry)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}

		if name == "" && password == "" && registry == "" {
			fmt.Println("Must specify either --name or --password or --registry")
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.PersistentFlags().StringVarP(&name, "name", "n", "", "Set npm account name")
	configCmd.PersistentFlags().StringVarP(&password, "pwd", "p", "", "Set npm account password")
	configCmd.PersistentFlags().StringVarP(&registry, "registry", "r", "", "Set npm registry")
}
