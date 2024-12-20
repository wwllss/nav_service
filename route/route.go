package route

import (
	"github.com/wwllss/zop"
)

func Register(z *zop.Zop) {
	nav := z.Group("/nav")
	nav.POST("", insertNavHandler)
	nav.POST("/notify", notifyNavHandler)
	nav.GET("/:os/:app/:version", queryNavHandler)
}
