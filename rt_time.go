package main

import (
	"github.com/wwllss/zop"
	"nav_service/hilog"
	"time"
)

var RtTime zop.NavHandlerFunc = func(c *zop.Context) {
	now := time.Now()
	c.Next()
	hilog.Infof("[%d]-route:%s, RT:%v", c.StatusCode, c.Path, time.Since(now))
}
