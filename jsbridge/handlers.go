package jsbridge

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/wwllss/zop"
	"html/template"
	"nav_service/dao"
	"nav_service/route"
	"net/http"
	"reflect"
	"sort"
	"strings"
)

var jsBridgeHandler zop.NavHandlerFunc = func(c *zop.Context) {
	jsBridge(c, dao.FindNewestVersion("Android", "com.ytmallapp"))
}

var jsBridgeByVersionHandler zop.NavHandlerFunc = func(c *zop.Context) {
	jsBridge(c, dao.FindVersion("Android", "com.ytmallapp", c.Param("v")))
}

func jsBridge(c *zop.Context, nv dao.Version) {
	common, android, ios := findCommon(getJsBridgeByVersion(nv, "cn.hipac.vm.webview.HvmJsBridge"),
		getJsBridgeByVersion(dao.FindVersion("iOS", "Mall", nv.Version), "JSBridgeProtocol"))
	mallCommon := make([]jsBridgeData, 0)
	for i := 0; i < len(common); {
		data := common[i]
		if data.IsCommon {
			i++
			continue
		}
		common = append(common[:i], common[i+1:]...)
		mallCommon = append(mallCommon, data)
	}
	for i := range common {
		common[i].Index = i + 1
	}
	for i := range mallCommon {
		mallCommon[i].Index = i + 1
	}
	for i := range android {
		android[i].Index = i + 1
	}
	for i := range ios {
		data := ios[i]
		data.Index = i + 1
		if len(data.Comment) == 0 {
			ios[i].Comment = firstLower(data.Token)
		}
	}
	data := jsHtml{
		Version: nv.Version,
		Common:  common,
		Mall: appBridge{
			Common:  mallCommon,
			Android: android,
			IOS:     ios,
		},
	}
	t := template.Must(template.ParseFiles("template/jsbridge.html"))
	if err := t.Execute(c.Writer, data); err != nil {
		c.String(http.StatusInternalServerError, "500 SERVER ERROR")
	}
}

func findCommon(a []jsBridgeData, i []jsBridgeData) (common []jsBridgeData, android []jsBridgeData, ios []jsBridgeData) {
	android = a
	ios = i
	am := make(map[string]jsBridgeData)
	for _, j := range android {
		am[j.Token] = j
	}
	cm := make(map[string]jsBridgeData)
	for _, j := range ios {
		aj, ok := am[firstLower(j.Token)]
		if !ok {
			continue
		}
		common = append(common, aj)
		cm[aj.Token] = aj
	}
	for i := 0; i < len(android); {
		_, ok := cm[android[i].Token]
		if !ok {
			i++
			continue
		}
		android = append(android[:i], android[i+1:]...)
	}
	for i := 0; i < len(ios); {
		_, ok := cm[firstLower(ios[i].Token)]
		if !ok {
			i++
			continue
		}
		ios = append(ios[:i], ios[i+1:]...)
	}
	return
}

func getJsBridgeByVersion(nv dao.Version, t string) []jsBridgeData {
	db := dao.GetDB()
	var navInfo route.NavInfo
	err := db.Last(&navInfo, "version_id=? and type=?", nv.ID, t).Error
	if err != nil {
		return nil
	}
	var service route.Service
	_ = json.Unmarshal([]byte(navInfo.Data), &service)
	return parseData(service)
}

func parseData(service route.Service) []jsBridgeData {
	jbdArr := make([]jsBridgeData, 0)
	for i, im := range service.ImplList {
		ca := strings.Split(im.Comment, "```")
		if len(ca) < 2 {
			ca = make([]string, 2)
		}
		jsBridge := jsBridgeData{
			Index:    i + 1,
			Token:    im.Token,
			Comment:  processComment(ca[0]),
			Example:  ca[1],
			IsCommon: strings.Contains(im.ImplName, "cn.hipac.vm.webview.HvmJsBridges"),
		}
		var str bytes.Buffer
		if err := json.Indent(&str, []byte(ca[1]), "", "    "); err != nil {
			jsBridge.Example = ca[1]
		} else {
			jsBridge.Example = str.String()
		}
		parseParam(&jsBridge)
		jbdArr = append(jbdArr, jsBridge)
	}
	return jbdArr
}

func processComment(c string) string {
	if strings.HasPrefix(c, " ") {
		return processComment(c[1:])
	}
	if !strings.HasPrefix(c, "@") {
		return c
	}
	split := strings.Split(c, " ")
	if len(split) == 0 {
		return c
	}
	return processComment(strings.Join(split[2:], " "))
}

func parseParam(js *jsBridgeData) {
	m := make(map[string]interface{})
	err := json.Unmarshal([]byte(js.Example), &m)
	if err != nil {
		a := make([]interface{}, 0)
		_ = json.Unmarshal([]byte(js.Example), &a)
		if len(a) == 0 {
			return
		}
		param := jsBridgeParam{
			Name: "入参",
		}
		parseArr(js, &param, "数组", a, 1)
		js.Params = append(js.Params, param)
	} else {
		if len(m) == 0 {
			return
		}
		parseObj(js, "入参", m)
	}
	for i, j := 0, len(js.Params)-1; i < j; i, j = i+1, j-1 {
		js.Params[i], js.Params[j] = js.Params[j], js.Params[i]
	}
}

func parseObj(j *jsBridgeData, name string, m map[string]interface{}) {
	ks := make([]string, 0)
	for k := range m {
		ks = append(ks, k)
	}
	sort.Strings(ks)
	param := jsBridgeParam{
		Name: name,
	}
	for _, k := range ks {
		v := m[k]
		vk := reflect.TypeOf(v).Kind()
		if vk == reflect.Slice {
			parseArr(j, &param, k, v.([]interface{}), 1)
		} else if vk == reflect.Map {
			param.Kvs = append(param.Kvs, KeyValue{Key: k, Value: "对象，详情见下个表格"})
			parseObj(j, k, v.(map[string]interface{}))
		} else {
			param.Kvs = append(param.Kvs, KeyValue{Key: k, Value: fmt.Sprintf("%v", v)})
		}
	}
	j.Params = append(j.Params, param)
}

func parseArr(j *jsBridgeData, param *jsBridgeParam, name string, a []interface{}, level int) {
	for _, v := range a {
		vk := reflect.TypeOf(v).Kind()
		if vk == reflect.Slice {
			parseArr(j, param, name, v.([]interface{}), level+1)
		} else if vk == reflect.Map {
			param.Kvs = append(param.Kvs, KeyValue{Key: name, Value: fmt.Sprintf("%d维对象数组，详情见下个表格", level)})
			parseObj(j, name, v.(map[string]interface{}))
		} else {
			param.Kvs = append(param.Kvs, KeyValue{Key: name, Value: fmt.Sprintf("%d维%s数组", level, vk.String())})
		}
	}
}

// FirstLower 字符串首字母小写
func firstLower(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToLower(s[:1]) + s[1:]
}
