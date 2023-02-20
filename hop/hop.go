package hop

import (
	"log"
	"net/http"
	"strings"
)

type NavHandlerFunc func(c *Context)

type NavGroup struct {
	hop         *Hop
	prefix      string
	parent      *NavGroup
	middlewares []NavHandlerFunc
}

type Hop struct {
	*NavGroup
	groups []*NavGroup
	nav    *nav
}

func New() *Hop {
	hop := &Hop{nav: newRouter()}
	hop.NavGroup = &NavGroup{hop: hop}
	hop.groups = []*NavGroup{hop.NavGroup}
	return hop
}

func (group *NavGroup) Group(prefix string) *NavGroup {
	hop := group.hop
	for _, g := range hop.groups {
		if g.prefix == prefix {
			return g
		}
	}
	newGroup := &NavGroup{
		hop:    hop,
		prefix: prefix,
		parent: group,
	}
	hop.groups = append(hop.groups, newGroup)
	return newGroup
}

func (group *NavGroup) Use(middlewares ...NavHandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

func (group *NavGroup) AddRoute(method string, path string, handler NavHandlerFunc) {
	log.Printf("AddRoute by Group : %s", group.prefix)
	group.hop.nav.AddRoute(method, group.prefix+path, handler)
}

func (group *NavGroup) GET(path string, handler NavHandlerFunc) {
	group.AddRoute("GET", path, handler)
}

func (group *NavGroup) POST(path string, handler NavHandlerFunc) {
	group.AddRoute("POST", path, handler)
}

func (hop *Hop) Run(addr string) (err error) {
	return http.ListenAndServe(addr, hop)
}

func (hop *Hop) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	middlewares := make([]NavHandlerFunc, 0)
	for _, group := range hop.groups {
		if strings.HasPrefix(request.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := NewContext(writer, request)
	c.handlers = middlewares
	hop.nav.Handle(c)
}
