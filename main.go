package main

import (
	"github.com/wwllss/zop"
	"nav_service/config"
	"nav_service/dao"
	"nav_service/hione"
	"nav_service/jsbridge"
	"nav_service/lottery"
	"nav_service/route"
)

func main() {
	db, _ := dao.GetDB().DB()
	defer func() {
		_ = db.Close()
	}()
	z := zop.New()
	z.Use(RtTime)
	z.Use(Recovery)
	route.Register(z)
	hione.Register(z)
	jsbridge.Register(z)
	lottery.Register(z)
	c := config.GetConfig()
	if err := z.Run(":" + c.Hop.Port); err != nil {
		panic(err)
	}
}
