package wstf

import (
	"regexp"
)

// @see https://github.com/yuya-takeyama/regexprouter/blob/master/regexprouter.go
// @see https://developer.mozilla.org/en-US/docs/Web/JavaScript/Guide/Regular_Expressions
// @see README.md#Parameter-Pattern
// > The **Parameter Pattern** conforms to the form: `{parameter-name:regexp-expression}`
// This regexp is used to find parameters and their corresponding configure
var patternRegexp = regexp.MustCompile(`\{([^\}:]+)(?::([^\}]+))?\}`)

type Route struct {
	// The original pattern given by initializer.
	Pattern string
	// The Regexp for matching path.
	PathRegexp *regexp.Regexp
	// Regexp for matching children.
	PrefixRegexp *regexp.Regexp
	// Parameter names in given pattern.
	ParamNames []string
	// Processors take effects only when Router is nil.
	Processors map[string]Processors
	// Router attached to the route.
	Router *Router
}

// Create a new route with a specific path.
func NewRoute(pattern string, router *Router) *Route {
	// Pre-process the pattern.
	//pattern = strings.Replace(pattern, "*", `.*`, -1)
	newPattern, err := generatePatternForRegexp(pattern)
	// Panic for abnormal patterns.
	if err != nil {
		panic("pattern [" + newPattern + "] is not recognized as normal.")
	}

	return &Route{
		pattern,
		regexp.MustCompile("^" + newPattern + "$"),
		regexp.MustCompile("^" + newPattern + "(.*)$"),
		findParamNames(pattern),
		make(map[string]Processors),
		router,
	}
}

// GivenPath > whether match this route.
func (m *Route) Handle(remainingPath string, req *Request, res *Response, next func()) {
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
