package lottery

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
	"nav_service/dao"
	"nav_service/hilog"
	"nav_service/hop"
	"nav_service/utils"
	"net/http"
	"reflect"
	"strings"
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

func todayNum(c *hop.Context, t Type) {
	rl := &randomLottery{
		Type: t,
	}
	err := dao.GetDB().Last(rl, "type = ?", t).Error
	if err != nil {
		c.String(http.StatusInternalServerError, "database error")
		return
	}
	if !utils.IsToday(rl.CreatedAt) {
		c.String(http.StatusOK, fmt.Sprintf("job not started, load num is:\n%s", rl.Nums))
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
