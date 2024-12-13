package route

import (
	"bytes"
	"encoding/json"
	"github.com/CatchZeng/dingtalk"
	"html/template"
	"nav_service/config"
	"nav_service/dao"
	"nav_service/hilog"
	"nav_service/hop"
	"net/http"
	"strings"
)

func init() {
	db := dao.GetDB()
	_ = db.AutoMigrate(&NavInfo{})
}

var insertNavHandler hop.NavHandlerFunc = func(c *hop.Context) {
	nv := dao.FindVersion(
		c.PostForm("os"),
		c.PostForm("app"),
		c.PostForm("version"),
	)
	db := dao.GetDB()
	thisInfo := &NavInfo{}
	lastInfo := &NavInfo{}
	formType := c.PostForm("type")
	formData := c.PostForm("data")
	thisErr := db.Last(thisInfo, "version_id=? and type=?", nv.ID, formType).Error
	lastErr := db.Last(lastInfo, "version_id=? and type=?", nv.LastVersionID, formType).Error

	if thisErr == nil {
		db.Model(thisInfo).Update("data", formData)
		thisInfo.Data = formData
	} else {
		thisInfo.VersionID = nv.ID
		thisInfo.Type = formType
		thisInfo.Data = formData
		db.Create(thisInfo)
	}
	if lastErr == nil {
		diffArr := diff(*lastInfo, *thisInfo)
		bytes, _ := json.Marshal(diffArr)
		thisInfo.DiffData = string(bytes)
		db.Model(thisInfo).Update("diff_data", thisInfo.DiffData)
	}
}

var notifyNavHandler hop.NavHandlerFunc = func(c *hop.Context) {
	nv := dao.FindVersion(
		c.PostForm("os"),
		c.PostForm("app"),
		c.PostForm("version"),
	)
	db := dao.GetDB()
	var infos []NavInfo
	err := db.Find(&infos, "version_id=?", nv.ID).Error
	if err != nil || len(infos) == 0 {
		c.String(http.StatusNotFound, "404 NOT FOUND")
		return
	}
	var diff []string
	for _, info := range infos {
		var temp []string
		_ = json.Unmarshal([]byte(info.DiffData), &temp)
		diff = append(diff, temp...)
	}
	if len(diff) == 0 {
		c.String(http.StatusOK, "Modify Nothing")
		return
	}
	ding(diff, nv)
}

func ding(arr []string, v dao.Version) {
	if len(arr) <= 0 {
		return
	}
	//TODO:token
	accessToken := "xxx"
	secret := "xxx"
	client := dingtalk.NewClient(accessToken, secret)
	t := template.Must(template.ParseFiles("template/nav_modify.md"))
	var sb strings.Builder
	hopConfig := config.GetConfig().Hop
	_ = t.Execute(&sb, mdData{v, arr[:5], hopConfig.Host, hopConfig.Port})
	message := dingtalk.NewMarkdownMessage().SetMarkdown(v.Os+v.Version+"路由修改", sb.String())
	_, _ = client.Send(message)
}

var queryNavHandler hop.NavHandlerFunc = func(c *hop.Context) {
	db := dao.GetDB()
	os := c.Param("os")
	app := c.Param("app")
	version := c.Param("version")
	nv := dao.FindVersion(os, app, version)
	var infos []NavInfo
	err := db.Find(&infos, "version_id=?", nv.ID).Error
	if err != nil || len(infos) == 0 {
		c.String(http.StatusNotFound, "404 NOT FOUND")
		return
	}
	data := htmlData{Version: nv}
	for _, info := range infos {
		if isActivity(info.Type) {
			_ = json.Unmarshal([]byte(info.Data), &data.Routes)
			_ = json.Unmarshal([]byte(info.DiffData), &data.RoutesDiff)
		} else {
			var serviceData htmlServiceData
			_ = json.Unmarshal([]byte(info.Data), &serviceData.Service)
			for i, im := range serviceData.Service.ImplList {
				a := strings.Split(im.Comment, "```")
				if len(a) != 3 {
					continue
				}
				var str bytes.Buffer
				if err := json.Indent(&str, []byte(a[1]), "", "    "); err != nil {
					continue
				}
				serviceData.Service.ImplList[i].ParamJson = str.String()
			}
			_ = json.Unmarshal([]byte(info.DiffData), &serviceData.ServiceDiff)
			data.Services = append(data.Services, serviceData)
		}
	}
	t := template.Must(template.ParseFiles("template/nav.html"))
	if err := t.Execute(c.Writer, data); err != nil {
		c.String(http.StatusInternalServerError, "500 SERVER ERROR")
	}
}

func diff(old NavInfo, new NavInfo) []string {
	if old.Data == new.Data {
		return nil
	}
	if isActivity(new.Type) {
		var oldData []route
		_ = json.Unmarshal([]byte(old.Data), &oldData)
		var co []Compare
		for _, o := range oldData {
			co = append(co, o)
		}
		var newData []route
		_ = json.Unmarshal([]byte(new.Data), &newData)
		var cn []Compare
		for _, n := range newData {
			cn = append(cn, n)
		}
		diff := Diff(co, cn)
		hilog.Infof("Diff %v", diff)
		return diff
	} else {
		var oldData Service
		_ = json.Unmarshal([]byte(old.Data), &oldData)
		var co []Compare
		for _, o := range oldData.ImplList {
			co = append(co, o)
		}
		var newData Service
		_ = json.Unmarshal([]byte(new.Data), &newData)
		var cn []Compare
		for _, n := range newData.ImplList {
			cn = append(cn, n)
		}
		diff := Diff(co, cn)
		hilog.Infof("Diff %v", diff)
		return diff
	}
}

func isActivity(t string) bool {
	return "activity" == t
}
