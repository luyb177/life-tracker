package summary

import (
	"context"
	"time"

	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/luyb177/life-tracker/backend/cmd/cron/internal"
	"github.com/luyb177/life-tracker/backend/internal/constvar"
	cronLogic "github.com/luyb177/life-tracker/backend/internal/logic/cron"
)

var periodType uint8

// Cmd cron summary 子命令
var Cmd = &cobra.Command{
	Use:   "summary",
	Short: "执行 AI 总结（日报/周报/月报/年报）",
	RunE:  run,
}

func init() {
	Cmd.Flags().Uint8VarP(&periodType, "type", "t", 1, "周期类型: 1=日报, 2=周报, 3=月报, 4=年报")
}

func run(_ *cobra.Command, _ []string) error {
	defer internal.SvcCtx.RedisClient.Client.Close()

	label := labelCN(periodType)
	logx.Infof("starting AI %s summary...", label)

	userIDs, err := internal.SvcCtx.Repos.User.ListIDs(context.Background())
	if err != nil {
		logx.Errorf("list users for AI %s summary failed: %v", label, err)
		return err
	}

	for _, userID := range userIDs {
		if err := cronLogic.Run(context.Background(), internal.SvcCtx, periodType, userID, time.Time{}); err != nil {
			logx.Errorf("AI %s summary failed for user %d: %v", label, userID, err)
			continue
		}
	}

	logx.Infof("AI %s summary completed", label)
	return nil
}

func labelCN(t uint8) string {
	switch t {
	case constvar.SummaryPeriodTypeDay:
		return "日报"
	case constvar.SummaryPeriodTypeWeek:
		return "周报"
	case constvar.SummaryPeriodTypeMonth:
		return "月报"
	case constvar.SummaryPeriodTypeYear:
		return "年报"
	default:
		return "总结"
	}
}
