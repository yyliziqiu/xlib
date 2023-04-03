package app

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var _rootCmd *cobra.Command

func ExecCmd(app *App, cmds ...func() *cobra.Command) {
	initRootCmd(app)

	for _, cmd := range cmds {
		_rootCmd.AddCommand(cmd())
	}

	err := _rootCmd.Execute()
	if err != nil {
		fmt.Printf("Exec cmd failed, error: %v.", err)
		os.Exit(1)
	}
}

func initRootCmd(app *App) {
	var (
		config = new(string)
		logdir = new(string)
	)

	_rootCmd = &cobra.Command{
		Version: app.Version,
		Use:     app.Name,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Use %s.exe -h or --help for help.\n", app.Name)
		},
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if strings.TrimSpace(*config) != "" {
				app.ConfigPath = *config
			}
			if strings.TrimSpace(*logdir) != "" {
				app.LogPath = *logdir
			}
			err := app.Init()
			if err != nil {
				fmt.Printf("Init app failed, error: %v\n", err)
				os.Exit(1)
			}
		},
	}

	_rootCmd.PersistentFlags().StringVarP(config, "config", "c", "", "config path")
	_rootCmd.PersistentFlags().StringVarP(logdir, "logdir", "d", "", "logdir path")
}
