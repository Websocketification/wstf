package wstf

// If the given path matches this route,
// it should exactly match the given pattern.
func (m *Route) MatchPath(remainingPath string, request *Request, response *Response) bool {
	//fmt.Println("Matching: ", m.Pattern, remainingPath, Stringify(request))
	// It can only have one match since the pattern is `^(.+)$` like.
	matches := m.PathRegexp.FindAllStringSubmatch(remainingPath, 1)
	if len(matches) == 0 {
		return false
	}
	if len(m.ParamNames) > 0 {
		for i := 0; i < len(m.ParamNames); i++ {
			request.Params[m.ParamNames[i]] = matches[0][i+1]
		}
	}
	return true
	matches2 := matches[0][1:]
	if len(matches2) > 0 {
		return true
	} else {
		return false
	}
	i := 0
	for _, match := range matches {
		if len(match) != 0 {
			return true
		} else {
			i++
		}
	}
	return false
}

// Whether pattern matches prefix path of given path.
// Return whether it matches and the remainingPath if it matches.
func (m *Route) MatchChildren(remainingPath string, request *Request, response *Response) (bool, string) {
	//fmt.Println("Matching PrefixPath: ", m.Pattern, remainingPath, Stringify(request))
	matches := m.PrefixRegexp.FindAllStringSubmatch(remainingPath, 1)
	if len(matches) == 0 {
		return false, ""
	}
	matches2 := matches[0][1:]
	if len(m.ParamNames) > 0 {
		for i := 0; i < len(m.ParamNames); i++ {
			request.Params[m.ParamNames[i]] = matches[0][i+1]
		}
	}
	return true, matches2[len(matches2)-1]
	i := 0
	for _, match := range matches {
		if len(match) != 0 {
			return true, matches2[len(matches2)-1]
		} else {
			i++
		}
	}
	return false, ""
}
