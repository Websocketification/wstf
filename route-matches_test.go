package wstf

import (
	"testing"
)

func TestRoute_Match(t *testing.T) {
	TestRoute_MatchFull(t)
	TestRoute_MatchParameterPattern(t)
	TestRoute_MatchAsteriskMark(t)
	TestRoute_MatchRegexp(t)
}

func TestRoute_MatchFull(t *testing.T) {
	testMatch(t, `/users`, "/users", true, true, "")
	testMatch(t, `/users`, "/user", false, false, "")
	testMatch(t, `/users`, "/users/fisher", false, true, "/fisher")
}

func TestRoute_MatchParameterPattern(t *testing.T) {
	testMatch(t, `/users/{userName}`, "/users/fisher", true, true, "")
	testMatch(t, `/users-{userName}`, "/users-fisher", true, true, "")
	testMatch(t, `/users/:userName`, "/users/fisher", false, false, "")
}

// Test the `*` asterisk mark.
func TestRoute_MatchAsteriskMark(t *testing.T) {
	testMatch(t, `/users/{userName}/.*`, "/users/fisher", false, false, "")
	testMatch(t, `/users/{userName}/.*`, "/users/fisher/Whatever", true, true, "")
	testMatch(t, `/users/{userName}/.*`, "/users/fisher/Whatever/etc", true, true, "")
	testMatch(t, `/users/{userName}/[^/]*`, "/users/fisher/Whatever/etc", false, true, "/etc")
}

func TestRoute_MatchRegexp(t *testing.T) {
	testMatch(t, `/users/{userName:[a-zA-Z0-9-_]+}`, "/users/zb-fisher", true, true, "")
	testMatch(t, `/users/{userName:[a-zA-Z0-9-_]+}`, "/users/zb.fisher", false, true, ".fisher")
	testMatch(t, `/users/{userName:[a-zA-Z0-9-_]+}`, "/users/zb/fisher", false, true, "/fisher")
	testMatch(t, `/files/{filename:[a-zA-Z0-9-_]+}\.{extension:[a-zA-Z0-9]+}`, "/files/caillou.mp3", true, true, "")
	//testMatch(t, `/files/{filename:[a-zA-Z0-9-_]+}\.{extension:[a-zA-Z0-9]{1,4}}`, "/files/caillou.mp3", true, true, "")
}

func TestRoute_MatchPanics(t *testing.T) {
	// This following tests should trigger panics.
	testMatch(t, `*/panic-now`, "/users/fisher", false, false, "")
	testMatch(t, `/wrong-regular-expression/panic-**-now`, "/users/fisher", false, false, "")
}

// Test Route#MatchPath() and Route#MatchChildren() methods.
func testMatch(t *testing.T, pattern, path string, expected bool, expectedMatchChildren bool, expectedRemainingPath string) {
	route := NewRoute(pattern, nil)
	req := NewFakeRequest()
	res := NewFakeResponse()
	matched := route.MatchPath(path, req, res)
	if matched != expected {
		t.Fatal("Failed to match path for [", pattern, path, "], expected [", expected, "] vs got [", matched, "]")
	}

	// Test matching children.
	matched, remainingPath := route.MatchChildren(path, req, res)
	if matched != expectedMatchChildren || expectedRemainingPath != remainingPath {
		t.Fatal("Failed to match children for [", pattern, path, "], expected [", expectedMatchChildren, expectedRemainingPath, "] vs got [", matched, remainingPath, "]")
	}
}
