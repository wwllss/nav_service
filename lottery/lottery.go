package lottery

import (
	"github.com/wwllss/zop"
)

type Type = string

const (
	SSQ Type = "SSQ"
	DLT Type = "DLT"
)

// lottery 用于存储多个区的配置参数
type lottery struct {
	Type       Type
	zones      []config // 多个区配置
	duplicates bool     // 是否允许重复号码
	sorted     bool     // 是否排序号码
}

// config 用于存储每个区的配置参数
type config struct {
	count int // 每个区的号码个数
	max   int // 每个区的号码范围
}

func Register(z *zop.Zop) {
	g := z.Group("/lottery")
	g.GET("/ssq", ssqHandler)
	g.GET("/dlt", dltHandler)
	randomEveryday()
}
