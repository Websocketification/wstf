package wstf

import (
	"net/http"
)

const (
	// General method that matches all method.
	MethodAll = "ALL"

	MethodLink = "LINK"

	MethodUnlink = "UNLINK"
)

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

func (m *Route) Head(processor Processor) *Route {
	return m.ListenMethod(http.MethodHead, processor)
}

func (m *Route) Options(processor Processor) *Route {
	return m.ListenMethod(http.MethodOptions, processor)
}

func (m *Route) Link(processor Processor) *Route {
	return m.ListenMethod(MethodLink, processor)
}

func (m *Route) Unlink(processor Processor) *Route {
	return m.ListenMethod(MethodUnlink, processor)
}

func (m *Route) TheEnd() {
	// Do Nothing to Explicitly End the Route.
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
