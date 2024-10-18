package monitor

import (
	"gitlab.com/cher8/lion/common/robot"
	"testing"
)

func TestFeishuRobot(t *testing.T) {
	const url = "https://open.feishu.cn/open-apis/bot/v2/hook/3e8eae6a-50ee-4b8d-b878-3de14bf2be27"
	fs := &robot.FeishuRobot{
		Url: url,
	}
	fs.Send2robot("hello cher8!")
}
