package wstf

import (
	"regexp"
	"strings"
	"fmt"
)

// @see https://github.com/yuya-takeyama/regexprouter/blob/master/regexprouter.go
// @see https://developer.mozilla.org/en-US/docs/Web/JavaScript/Guide/Regular_Expressions
// @see README.md#Parameter-Pattern
// > The **Parameter Pattern** conforms to the form: `{parameter-name:regexp-expression}`
// This regexp is used to find parameters and their corresponding configure
var patternRegexp = regexp.MustCompile(`\{([^\}:]+)(?::([^\}]+))?\}`)

type Route struct {
	// The Regexp for matching path.
	Regexp *regexp.Regexp
	// Regexp for matching children.
	RegexpPrefix *regexp.Regexp
	// The original pattern given by initializer.
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

	// Pre-process the pattern.
	pattern = strings.Replace(pattern, "*", `[^/]*`, -1)

	// Find all the parameters.
	parameters := patternRegexp.FindAllStringSubmatch(pattern, -1)
	for _, parameter := range parameters {
		name := parameter[1]
		m.ParamNames = append(m.ParamNames, name)
	}

	// Set up regexps.
	newPattern := patternRegexp.ReplaceAllStringFunc(pattern, func(subMatch string) string {
		// The sub match string conforms the parameter pattern: `{parameter-name:regexp-expression}`.
		foos := strings.SplitN(subMatch, ":", 2)
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
	//m.ListenMethod(MethodAll, processors...)
	return m
}

// GivenPath > whether match this route.
func (m Route) Handle(remainingPath string, req Request, res Response, next func()) {
	if m.MatchPath(remainingPath, req, res) {
		if m.Router != nil {
			m.Router.Handle("", req, res, next)
		} else {
			HandleMethods(m.Processors[MethodAll], req, res, func() {
				HandleMethods(m.Processors[req.Method], req, res, next)
			})
		}
	} else if matched, remainingPath := m.MatchChildren(remainingPath, req, res); m.Router != nil && matched {
		// Here can be optimized!
		m.Router.Handle(remainingPath, req, res, next)
	} else {
		next()
	}
}
