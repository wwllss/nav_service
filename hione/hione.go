package hione

import (
	"gorm.io/gorm"
	"nav_service/dao"
	"nav_service/hop"
	"net/http"
)

type apkTestUpdate struct {
	*gorm.Model
	Version          string
	HioneProjectCode string
	Builder          string
	DownloadUrl      string
	Time             string
	/*BuildStatus      string
	FailType         string
	FailDesc         string
	VersionCode      uint
	ProjectName      string
	DepConfigs       string
	DebugEnable      string
	Platform         string
	DownloadQRCode   string `gorm:"column:downloadQRCode"`
	UpdateDesc       string `gorm:"column:updateDesc"`
	BizContent       string `gorm:"column:bizContent"`
	BuildEnv         string
	BuildTime        string*/
}

func (abi apkTestUpdate) TableName() string {
	return "apk_test_update"
}

func init() {
	_ = dao.GetDB().AutoMigrate(&apkTestUpdate{})
}

func Register(h *hop.Hop) {
	g := h.Group("/hione")
	g.POST("/apk", func(c *hop.Context) {
		db := dao.GetDB()
		info := &apkTestUpdate{
			Version:          c.PostForm("v"),
			HioneProjectCode: c.PostForm("c"),
			Builder:          c.PostForm("b"),
			DownloadUrl:      c.PostForm("d"),
			Time:             c.PostForm("t"),
		}
		if info.DownloadUrl == "" {
			c.String(http.StatusNotFound, "insert failed:DownloadUrl must not be nil")
			return
		}
		if info.Time == "" {
			c.String(http.StatusNotFound, "insert failed:Time must not be nil")
			return
		}
		if err := db.Create(info).Error; err != nil {
			c.String(http.StatusNotFound, "insert failed")
		} else {
			c.String(http.StatusOK, "insert success")
		}
	})
	g.POST("/update", func(c *hop.Context) {
		db := dao.GetDB()
		info := &apkTestUpdate{}
		if err := db.Last(
			info,
			"hione_project_code=?",
			c.PostForm("hione"),
		).Error; err != nil {
			c.Json(http.StatusNotFound, map[string]interface{}{
				"msg": "未找到",
			})
			return
		}
		if info.Time == c.PostForm("time") {
			c.Json(http.StatusNotFound, map[string]interface{}{
				"msg": "无需更新",
			})
			return
		}
		c.Json(http.StatusOK, map[string]interface{}{
			"version": info.Version,
			"time":    info.Time,
			"url":     info.DownloadUrl,
			"builder": info.Builder,
		})
	})
}
