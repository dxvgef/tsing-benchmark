package tsing_gateway

import (
	"net/http"
)

type Engine struct {
	hosts map[string]map[string]map[string]RouteNode
}

func New() *Engine {
	var proxy Engine
	proxy.hosts = make(map[string]map[string]map[string]RouteNode)
	return &proxy
}

func (p *Engine) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	_, status := p.matchRoute(req)
	if status != http.StatusOK {
		resp.WriteHeader(status)
		return
	}
	resp.WriteHeader(204)
}
