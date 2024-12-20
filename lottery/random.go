package lottery

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"github.com/wwllss/zop"
	"gorm.io/gorm"
	"math/rand"
	"nav_service/dao"
	"nav_service/hilog"
	"nav_service/utils"
	"net/http"
	"reflect"
	"sort"
	"strings"
	"time"
)

type randomLottery struct {
	gorm.Model
	Type Type
	Nums string
}

func randomEveryday() {
	err := dao.GetDB().AutoMigrate(&randomLottery{})
	if err != nil {
		return
	}
	c := cron.New()
	_, _ = c.AddFunc(ssqTime, func() {
		job(ssqConfig)
	})
	_, _ = c.AddFunc(dltTime, func() {
		job(dltConfig)
	})
	c.Start()
}

func todayNum(c *zop.Context, t Type) {
	rl := &randomLottery{
		Type: t,
	}
	err := dao.GetDB().Last(rl, "type = ?", t).Error
	if err != nil {
		c.String(http.StatusInternalServerError, "未查到")
		return
	}
	if !utils.IsToday(rl.CreatedAt) {
		c.String(http.StatusOK, fmt.Sprintf("任务未开始，上期号码为:\n%s", rl.Nums))
		return
	}
	c.String(http.StatusOK, rl.Nums)
}

func job(l lottery) {
	hilog.Infof("开始选票 --- %s", l.Type)
	tickets := make([][][]int, 0)
	for {
		if len(tickets) == 5 {
			break
		}
		t := findLucky(l)
		if t != nil {
			tickets = append(tickets, t)
		}
	}
	st := format(tickets)
	hilog.Infof("今日票选如下：\n%s", st)
	dao.GetDB().Create(&randomLottery{
		Type: l.Type,
		Nums: st,
	})
}

func format(tickets [][][]int) string {
	l := len(tickets)
	if l <= 0 {
		return ""
	}
	ss := make([]string, l)
	for i, ticket := range tickets {
		ss[i] = fmt.Sprintf("%v", ticket)
	}
	return strings.Join(ss, "\n")
}

func findLucky(l lottery) [][]int {
	lucky := generateLotteryNumbers(l)
	tickets := make([][][]int, 5)
	for i := range tickets {
		tickets[i] = generateLotteryNumbers(l)
	}
	for _, ticket := range tickets {
		if reflect.DeepEqual(ticket, lucky) {
			hilog.Infof("一等奖：%v", ticket)
			return ticket
		}
	}
	return nil
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
