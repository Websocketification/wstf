package wstf

import (
	"net/http"
)

// General method that matches all method.
const MethodAll = "ALL"

func (m *Route) All(processor Processor) *Route {
	return m.ListenMethod(MethodAll, processor)
}

func (m *Route) Get(processor Processor) *Route {
	return m.ListenMethod(http.MethodGet, processor)
}

func (m *Route) Post(processor Processor) *Route {
	return m.ListenMethod(http.MethodPost, processor)
}

func (m *Route) Patch(processor Processor) *Route {
	return m.ListenMethod(http.MethodPatch, processor)
}

func (m *Route) Put(processor Processor) *Route {
	return m.ListenMethod(http.MethodPut, processor)
}

func (m *Route) Delete(processor Processor) *Route {
	return m.ListenMethod(http.MethodDelete, processor)
}

// Add processors to listen specific method.
func (m *Route) ListenMethod(method string, processor Processor) *Route {
	if m.Processors[method] == nil {
		m.Processors[method] = Processors{processor}
	} else {
		m.Processors[method] = append(m.Processors[method], processor)
	}
	return m
}
