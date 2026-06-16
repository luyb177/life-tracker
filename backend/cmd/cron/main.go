package main

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/core/conf"

	"github.com/luyb177/life-tracker/backend/cmd/cron/internal"
	"github.com/luyb177/life-tracker/backend/cmd/cron/summary"
	"github.com/luyb177/life-tracker/backend/internal/config"
	"github.com/luyb177/life-tracker/backend/internal/svc"
)

const codeFailure = 1

var (
	confPath string

	rootCmd = &cobra.Command{
		Use:   "cron",
		Short: "life-tracker cron jobs",
		Long:  "life-tracker 定时任务调度",
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(codeFailure)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&confPath, "config", "etc/life_tracker.yaml", "config file")

	// 注册子命令
	rootCmd.AddCommand(summary.Cmd)
}

func initConfig() {
	var c config.Config
	conf.MustLoad(confPath, &c)
	internal.SvcCtx = svc.NewServiceContext(c)
}

func main() {
	Execute()
}
