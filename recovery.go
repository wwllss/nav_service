package main

import (
	"fmt"
	"github.com/wwllss/zop"
	"nav_service/hilog"
	"net/http"
	"runtime"
	"strings"
)

var Recovery zop.NavHandlerFunc = func(c *zop.Context) {
	defer func() {
		if err := recover(); err != nil {
			message := fmt.Sprintf("%s", err)
			hilog.Infof("%s\n\n", trace(message))
			c.String(http.StatusInternalServerError, "Internal Server Error")
		}
	}()
	c.Next()
}

func trace(message string) string {
	var pcs [32]uintptr
	n := runtime.Callers(3, pcs[:]) // skip first 3 caller

	var str strings.Builder
	str.WriteString(message + "\nTraceback:")
	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return str.String()
}
