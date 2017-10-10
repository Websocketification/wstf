package wstf

import (
	"regexp"
	"strings"
	"fmt"
)

// @see https://github.com/yuya-takeyama/regexprouter/blob/master/regexprouter.go
var patternRegexp = regexp.MustCompile(`\{([^\}:]+)(?::([^\}]+))?\}`)

type Route struct {
	// Regexp for matching path.
	Regexp *regexp.Regexp
	// Regexp for matching children.
	RegexpPrefix *regexp.Regexp
	// Given pattern.
	Pattern string
	// Parameter names in given pattern.
	ParamNames []string
	// Processors take effects only when Router is nil.
	Processors map[string][]func(req Request, res Response, next func())
	Router     *Router
}

// Create a new route with a specific path.
func NewRoute(pattern string, router *Router) Route {
	m := Route{
		Pattern:    pattern,
		Processors: make(map[string][]func(req Request, res Response, next func())),
		Router:     router,
	}
	parameters := patternRegexp.FindAllStringSubmatch(pattern, -1)
	for _, parameter := range parameters {
		name := parameter[1]
		m.ParamNames = append(m.ParamNames, name)
	}
	newPattern := patternRegexp.ReplaceAllStringFunc(pattern, func(s string) string {
		foos := strings.SplitN(s, ":", 2)
		if len(foos) < 2 {
			return `([^/]+)`
		} else {
			return "(" + foos[1][0:len(foos[1])-1] + ")"
		}
	})
	fmt.Println("newPattern: ", newPattern, m.ParamNames)
	m.Regexp = regexp.MustCompile("^" + newPattern + "$")
	m.RegexpPrefix = regexp.MustCompile("^" + newPattern + "(.*)$")
	//for i, _ := range m.ParamNames {
	//	if i == 0 {
	//processors = append(processors, fn)
	//} else {
	//processors = append(processors, nil)
	//}
	//}

	//r.names = append(r.names, names)
	//m.ListenMethod(METHOD_ALL, processors...)
	return m
}

// GivenPath > whether match this route.
func (m Route) Handle(remainingPath string, req Request, res Response, next func()) {
	if m.Match(remainingPath, req, res) {
		if m.Router != nil {
			m.Router.Handle("", req, res, next)
		} else {
			HandleMethods(m.Processors[METHOD_ALL], req, res, func() {
				HandleMethods(m.Processors[req.Method], req, res, next)
			})
		}
	} else if matched, remainingPath := m.MatchPrefixPath(remainingPath, req, res); m.Router != nil && matched {
		// Here can be optimized!
		m.Router.Handle(remainingPath, req, res, next)
	} else {
		next()
	}
}
