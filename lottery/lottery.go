package lottery

import (
	"math/rand"
	"nav_service/hop"
	"sort"
	"time"
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

// generateLotteryNumbers 生成多个区的号码
func generateLotteryNumbers(config lottery) [][]int {
	// 使用当前时间戳创建一个新的随机源
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	var fullNums [][]int
	// 生成每个区的号码
	for _, zoneConfig := range config.zones {
		nums := generateNumbers(r, zoneConfig.count, zoneConfig.max, config.duplicates)
		fullNums = append(fullNums, nums)
	}

	// 如果需要排序，根据策略决定
	if config.sorted {
		for i := range fullNums {
			sort.Ints(fullNums[i])
		}
	}

	return fullNums
}

// generateNumbers 根据配置生成号码
func generateNumbers(r *rand.Rand, count, numberRange int, allowDuplicates bool) []int {
	var numbers []int
	for len(numbers) < count {
		num := r.Intn(numberRange) + 1
		if allowDuplicates || !contains(numbers, num) {
			numbers = append(numbers, num)
		}
	}
	return numbers
}

// contains 检查切片中是否包含某个号码
func contains(slice []int, num int) bool {
	for _, n := range slice {
		if n == num {
			return true
		}
	}
	return false
}

func Register(h *hop.Hop) {
	g := h.Group("/lottery")
	g.GET("/ssq", ssqHandler)
	g.GET("/dlt", dltHandler)
	randomEveryday()
}
