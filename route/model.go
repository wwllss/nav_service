package route

import (
	"fmt"
	"gorm.io/gorm"
	"nav_service/dao"
)

type NavInfo struct {
	gorm.Model
	VersionID uint
	Type      string
	Data      string `gorm:"type:text"`
	DiffData  string `gorm:"type:text"`
}

type mdData struct {
	Version dao.Version
	Diff    []string
	Host    string
	Port    string
}

type htmlData struct {
	Version    dao.Version
	Routes     []route
	RoutesDiff []string
	Services   []htmlServiceData
}

type htmlServiceData struct {
	Service     Service
	ServiceDiff []string
}

type routeArg struct {
	ParamName string `json:"paramName"`
	FieldName string `json:"fieldName"`
}

func (ra routeArg) Key() string {
	return ra.ParamName
}

func (ra routeArg) Compare(o Compare) []string {
	ra2 := o.(routeArg)
	diff := make([]string, 0)
	if ra.ParamName != ra2.ParamName {
		diff = append(diff, fmt.Sprintf("    Modify:ParamName:%s -> %s", ra.ParamName, ra2.ParamName))
	}
	if ra.FieldName != ra2.FieldName {
		diff = append(diff, fmt.Sprintf("    Modify:FieldName:%s -> %s", ra.FieldName, ra2.FieldName))
	}
	return diff
}

type interceptor string

func (i interceptor) Key() string {
	return string(i)
}

func (i interceptor) Compare(Compare) []string {
	return make([]string, 0)
}

type route struct {
	Route           string        `json:"route"`
	ClassName       string        `json:"className"`
	Comment         string        `json:"comment"`
	ArgList         []routeArg    `json:"argList"`
	InterceptorList []interceptor `json:"interceptorList"`
}

func (r route) Key() string {
	return r.Route
}

func (r route) Compare(o Compare) []string {
	r2 := o.(route)
	diff := make([]string, 0)
	if r.ClassName != r2.ClassName {
		diff = append(diff, fmt.Sprintf("    Modify:ClassName:%s -> %s", r.ClassName, r2.ClassName))
	}
	if r.Comment != r2.Comment {
		diff = append(diff, fmt.Sprintf("    Modify:Comment:%s -> %s", r.Comment, r2.Comment))
	}
	//cast arg
	c1 := make([]Compare, 0)
	for _, arg := range r.ArgList {
		c1 = append(c1, arg)
	}
	c2 := make([]Compare, 0)
	for _, arg := range r2.ArgList {
		c2 = append(c2, arg)
	}
	argDiff := Diff(c1, c2)
	for i, d := range argDiff {
		argDiff[i] = "    Arg:" + d
	}
	//cast interceptor
	ci1 := make([]Compare, 0)
	for _, i := range r.InterceptorList {
		ci1 = append(ci1, i)
	}
	ci2 := make([]Compare, 0)
	for _, i := range r2.InterceptorList {
		ci2 = append(ci2, i)
	}
	interceptorDiff := Diff(ci1, ci2)
	for i, d := range interceptorDiff {
		interceptorDiff[i] = "    Interceptor:" + d
	}
	//merge
	diff = append(diff, append(argDiff, interceptorDiff...)...)
	if len(diff) > 0 {
		diff = append([]string{fmt.Sprintf("Modify:%s", r.Route)}, diff...)
	}
	return diff
}

type ServiceImpl struct {
	Token     string `json:"token"`
	ImplName  string `json:"implName"`
	Comment   string `json:"comment"`
	ParamJson string
}

func (si ServiceImpl) Key() string {
	return si.Token
}

func (si ServiceImpl) Compare(o Compare) []string {
	si2 := o.(ServiceImpl)
	diff := make([]string, 0)
	if si.ImplName != si2.ImplName {
		diff = append(diff, fmt.Sprintf("    Modify:ImplName:%s -> %s", si.ImplName, si2.ImplName))
	}
	if si.Comment != si2.Comment {
		diff = append(diff, fmt.Sprintf("    Modify:Comment:%s -> %s", si.Comment, si2.Comment))
	}
	return diff
}

type Service struct {
	Service  string        `json:"service"`
	ImplList []ServiceImpl `json:"implList"`
}
