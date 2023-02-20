package main

import (
	"log"
	"nav_service/hop"
	"time"
)

var RtTime hop.NavHandlerFunc = func(c *hop.Context) {
	now := time.Now()
	c.Next()
	log.Printf("[%d]-route:%s, RT:%v", c.StatusCode, c.Path, time.Since(now))
}
