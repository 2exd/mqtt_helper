package cmd

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"mqtt_helper/conf"
	"mqtt_helper/constants"
	"mqtt_helper/log"
	"mqtt_helper/structs"
	"os/signal"
	"syscall"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "start as server",
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.SetConfigName(constants.SERVER)
		viper.SetConfigType(constants.YAML)
		viper.AddConfigPath(".")
		viper.AddConfigPath("./conf/") // 配置文件的路径

		err := viper.ReadInConfig() // 找到并加载配置文件
		if err != nil {             // 处理错误
			panic(fmt.Errorf("Fatal error config file: %w \n", err))
		}

		config := conf.GetConfig()

		// 初始化 mqtt client
		_, err = structs.NewMqttServer(
			structs.PubTopics(config.PubTopics),
			structs.SubTopics(config.SubTopics),
		)
		if err != nil {
			log.Logger.Fatalf("create mqtt client failed, %v", err)
			return
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		var ctx, stop = signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		defer stop()

		ms := structs.GetMqttServerInstance()
		if err := ms.Run(ctx); err != nil {
			log.Logger.Errorf("start failed, %v", err)
		} else {
			ms.PublishServerDown()
			log.Logger.Info("app stop")
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
