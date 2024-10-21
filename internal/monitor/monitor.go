package monitor

import (
	"context"
	"fmt"
	"grs/internal/common/robot"
	"grs/internal/common/util"
	"grs/internal/dump"
	"grs/internal/index"
	"os"
	"time"
)

func DeletedUserViewKeys(ctx context.Context, times int, args ...string) {
	// 定时删除用户曝光列表
	ticker := time.NewTicker(time.Second * 86400 * time.Duration(times))
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			dump.DeletedStringUserViewKeys(ctx, args[0])
		}
	}
}

func UpdateIndex(forwardColsStr string, invertedColsStr string, times int, args ...string) {
	// 定时更新索引
	ticker := time.NewTicker(time.Second * time.Duration(times))
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			index.NewIndex("update", forwardColsStr, invertedColsStr, args...)
		}
	}
}

func CheckServerHealth(feishuUrl string, times int) {
	ticker := time.NewTicker(time.Second * time.Duration(times))
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			fs := &robot.FeishuRobot{Url: feishuUrl}
			fs.Send2robot(
				fmt.Sprintf("%s, %s, Everything is ok now!", util.GetNowTimeFmt(), os.Getenv("REC_ENV")))
		}
	}
}
