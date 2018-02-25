package wstf

import (
	"regexp"
	"strings"
)

// Find all the parameters in user defined pattern.
func findParamNames(pattern string) []string {
	var paramNames []string
	parameters := patternRegexp.FindAllStringSubmatch(pattern, -1)
	for _, parameter := range parameters {
		name := parameter[1]
		paramNames = append(paramNames, name)
	}
	return paramNames
}

// Set up pattern for regexp.
func generatePatternForRegexp(pattern string) (string, error) {
	pattern = patternRegexp.ReplaceAllStringFunc(pattern, func(subMatch string) string {
		// The sub match string conforms the parameter pattern: `{parameter-name:regexp-expression}`.
		foos := strings.SplitN(subMatch, ":", 2)
		if len(foos) < 2 {
			return `([^/]+)`
		} else {
			return "(" + foos[1][0:len(foos[1])-1] + ")"
		}
	})
	// Checking for abnormal patterns.
	_, err := regexp.Compile(pattern)
	return pattern, err
}
