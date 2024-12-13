package main

import (
	"nav_service/hilog"
	"nav_service/hop"
	"time"
)

var RtTime hop.NavHandlerFunc = func(c *hop.Context) {
	now := time.Now()
	c.Next()
	hilog.Infof("[%d]-route:%s, RT:%v", c.StatusCode, c.Path, time.Since(now))
}
