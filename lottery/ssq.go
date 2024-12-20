package lottery

import (
	"github.com/wwllss/zop"
)

var ssqTime = "0 18 * * 0,2,4"

var ssqConfig = lottery{
	Type: SSQ,
	zones: []config{
		{count: 6, max: 33}, // 第1区，6个号码，范围1-33
		{count: 1, max: 16}, // 第2区，1个号码，范围1-16
	},
	duplicates: false, // 不允许重复号码
	sorted:     true,  // 排序号码
}

var ssqHandler = func(c *zop.Context) {
	todayNum(c, SSQ)
}
