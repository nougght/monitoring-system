package cli

import (
	"agent/internal/model"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var rootArgs = struct {
	setupConfig model.ArgInfo
}{
	setupConfig: model.ArgInfo{
		Name:         "setupconfig",
		DefaultValue: "",
		Description:  "agent setup config path",
		Required:     true,
	},
}

type CLI struct {
	rootCmd         *cobra.Command
	setupConfigPath *string
}

func New() *CLI {
	cli := &CLI{}
	cli.rootCmd = cli.newRootCmd()
	cli.rootCmd.AddCommand(newRunCmd(cli.setupConfigPath))
	return cli
}

func (c *CLI) Execute() int {
	if err := c.rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}

func (c *CLI) newRootCmd() *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   "agent",
		Short: "Monitoring agent",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if *c.setupConfigPath == rootArgs.setupConfig.DefaultValue {
				log.Println("setup config flag not set")
			}
			// TODO: auto resolve config path
			return nil
		},
	}
	rootCmd.PersistentFlags().StringVar(c.setupConfigPath,
		rootArgs.setupConfig.Name,
		rootArgs.setupConfig.DefaultValue,
		rootArgs.setupConfig.Description)
	return rootCmd
}
