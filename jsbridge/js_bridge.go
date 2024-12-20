package jsbridge

import (
	"github.com/wwllss/zop"
)

func Register(z *zop.Zop) {
	jsBridge := z.Group("/jsbridge")
	jsBridge.GET("", jsBridgeHandler)
	jsBridge.GET("/:v", jsBridgeByVersionHandler)
}
