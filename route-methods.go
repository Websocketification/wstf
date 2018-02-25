package wstf

import (
	"net/http"
)

// General method that matches all method.
const MethodAll = "ALL"

func (m *Route) All(processors Processor) *Route {
	return m.ListenMethod(MethodAll, processors)
}

func (m *Route) Get(processors Processor) *Route {
	return m.ListenMethod(http.MethodGet, processors)
}

func (m *Route) Post(processors Processor) *Route {
	return m.ListenMethod(http.MethodPost, processors)
}

func (m *Route) Patch(processors Processor) *Route {
	return m.ListenMethod(http.MethodPatch, processors)
}

func (m *Route) Put(processors Processor) *Route {
	return m.ListenMethod(http.MethodPut, processors)
}

func (m *Route) Delete(processors Processor) *Route {
	return m.ListenMethod(http.MethodDelete, processors)
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
