package tsing_gateway

import (
	"net/http"
	"path"
	"strings"
	"unsafe"
)

type RouteNode struct {
	Service string `json:"service"`
}

func (p *Engine) PutRoute(hostname, path, method string, routeNode RouteNode) {
	if hostname == "" {
		hostname = "*"
	}
	if path == "" {
		path = "/"
	}
	if method == "" {
		method = "*"
	} else {
		method = strings.ToUpper(method)
	}
	if _, exist := p.hosts[hostname]; !exist {
		p.hosts[hostname] = make(map[string]map[string]RouteNode)
	}
	if _, exist := p.hosts[hostname][path]; !exist {
		p.hosts[hostname][path] = make(map[string]RouteNode)
	}
	p.hosts[hostname][path][method] = routeNode
}

func (p *Engine) matchRoute(req *http.Request) (routeNode RouteNode, status int) {
	path := req.URL.Path
	method := req.Method
	host := req.Host
	result := false

	// 匹配主机
	host, result = p.matchHost(host)
	if !result {
		status = http.StatusServiceUnavailable
		return
	}
	path, result = p.matchPath(host, path)
	if !result {
		status = http.StatusNotFound
		return
	}
	method, result = p.matchMethod(host, path, method)
	if !result {
		status = http.StatusMethodNotAllowed
		return
	}

	routeNode = p.hosts[host][path][method]
	status = http.StatusOK
	return
}

func (p *Engine) matchHost(reqHost string) (string, bool) {
	pos := strings.LastIndex(reqHost, ":")
	if pos > -1 {
		reqHost = reqHost[:pos]
	}
	if _, exist := p.hosts[reqHost]; exist {
		return reqHost, true
	}
	reqHost = "*"
	if _, exist := p.hosts[reqHost]; exist {
		return reqHost, true
	}
	return reqHost, false
}

func (p Engine) matchPath(reqHost, reqPath string) (string, bool) {
	if reqPath == "" {
		reqPath = "/"
	}
	if _, exist := p.hosts[reqHost][reqPath]; exist {
		return reqPath, true
	}
	pathLastChar := reqPath[len(reqPath)-1]
	if pathLastChar != 47 {
		pos := strings.LastIndex(reqPath, path.Base(reqPath))
		reqPath = reqPath[:pos]
	}
	reqPath = reqPath + "*"
	if _, exist := p.hosts[reqHost][reqPath]; exist {
		return reqPath, true
	}
	return reqPath, false
}

func (p *Engine) matchMethod(reqHost, reqPath, reqMethod string) (string, bool) {
	if _, exist := p.hosts[reqHost][reqPath][reqMethod]; exist {
		return reqMethod, true
	}
	reqMethod = "*"
	if _, exist := p.hosts[reqHost][reqPath][reqMethod]; exist {
		return reqMethod, true
	}
	return reqMethod, false
}

func strToBytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s)) // nolint
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h)) // nolint
}
