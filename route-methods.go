package wstf

import (
	"net/http"
)

// General method that matches all method.
const METHOD_ALL = "ALL"

func (m Route) All(processors func(req Request, res Response, next func())) Route {
	m.ListenMethod(METHOD_ALL, processors)
	return m
}

func (m Route) Get(processors func(req Request, res Response, next func())) Route {
	m.ListenMethod(http.MethodGet, processors)
	return m
}

func (m Route) Post(processors func(req Request, res Response, next func())) Route {
	m.ListenMethod(http.MethodPost, processors)
	return m
}

func (m Route) Patch(processors func(req Request, res Response, next func())) Route {
	m.ListenMethod(http.MethodPatch, processors)
	return m
}

func (m Route) Put(processors func(req Request, res Response, next func())) Route {
	m.ListenMethod(http.MethodPut, processors)
	return m
}

func (m Route) Delete(processors func(req Request, res Response, next func())) Route {
	m.ListenMethod(http.MethodDelete, processors)
	return m
}

// Add processors to listen specific method.
func (m Route) ListenMethod(method string, processor func(req Request, res Response, next func())) {
	if m.Processors[method] == nil {
		m.Processors[method] = []func(req Request, res Response, next func()){processor}
	} else {
		m.Processors[method] = append(m.Processors[method], processor)
	}
}
