package jsbridge

import (
	"nav_service/hop"
)

func Register(h *hop.Hop) {
	jsBridge := h.Group("/jsbridge")
	jsBridge.GET("", jsBridgeHandler)
	jsBridge.GET("/:v", jsBridgeByVersionHandler)
}
