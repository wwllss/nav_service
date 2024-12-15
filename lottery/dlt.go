package lottery

import (
	"nav_service/hop"
)

var dltTime = "0 18 * * 1,3,6"

var dltConfig = lottery{
	Type: DLT,
	zones: []config{
		{count: 5, max: 35}, // 第1区，5个号码，范围1-35
		{count: 2, max: 12}, // 第2区，1个号码，范围1-12
	},
	duplicates: false, // 不允许重复号码
	sorted:     true,  // 排序号码
}

var dltHandler = func(c *hop.Context) {
	todayNum(c, DLT)
}
