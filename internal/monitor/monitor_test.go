package monitor

import (
	"grs/internal/common/robot"
	"testing"
)

func TestFeishuRobot(t *testing.T) {
	const url = "https://open.feishu.cn/open-apis/bot/v2/hook/..."
	fs := &robot.FeishuRobot{
		Url: url,
	}
	fs.Send2robot("hello world!")
}
