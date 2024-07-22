package cmd

import (
	"fmt"
	"np/util"
	"os"

	"github.com/spf13/cobra"
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
		config, err := util.ReadConfig()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if name != "" {
			config["username"] = name
		}
		if password != "" {
			config["password"] = password
		}
		if registry != "" {
			config["registry"] = registry
		}
		if name == "" && password == "" && registry == "" {
			fmt.Println("Must specify either --name or --password or --registry")
		}
		err = util.WriteConfig(config)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.PersistentFlags().StringVarP(&name, "name", "n", "", "Set npm account name")
	configCmd.PersistentFlags().StringVarP(&password, "pwd", "p", "", "Set npm account password")
	configCmd.PersistentFlags().StringVarP(&registry, "registry", "r", "", "Set npm registry")
}
