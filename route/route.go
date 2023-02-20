package route

import (
	"nav_service/hop"
)

func Register(h *hop.Hop) {
	nav := h.Group("/nav")
	nav.POST("", insertNavHandler)
	nav.POST("/notify", notifyNavHandler)
	nav.GET("/:os/:app/:version", queryNavHandler)
}
