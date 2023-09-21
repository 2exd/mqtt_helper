package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"mqtt_helper/log"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "helper",
	Short: "MQTT helper",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// InitLogger
		log.InitLogger()
		log.Logger.Info("app start")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
}
