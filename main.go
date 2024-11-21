package main

import (
	"nav_service/config"
	"nav_service/dao"
	"nav_service/hione"
	"nav_service/hop"
	"nav_service/jsbridge"
	"nav_service/route"
)

func main() {
	db, _ := dao.GetDB().DB()
	defer func() {
		_ = db.Close()
	}()
	h := hop.New()
	h.Use(RtTime)
	h.Use(Recovery)
	route.Register(h)
	hione.Register(h)
	jsbridge.Register(h)
	c := config.GetConfig()
	if err := h.Run(":" + c.Hop.Port); err != nil {
		panic(err)
	}
}
