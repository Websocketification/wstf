package wstf

import (
	"testing"
)

func TestRoute_Match(t *testing.T) {
	TestRoute_MatchFull(t)
	TestRoute_MatchParameterPattern(t)
}

func TestRoute_MatchFull(t *testing.T) {
	testMatch(t, `/users`, "/users", true)
	testMatch(t, `/users`, "/user", false)
}

func TestRoute_MatchParameterPattern(t *testing.T) {
	testMatch(t, `/users/{userName}`, "/users/fisher", true)
	testMatch(t, `/users-{userName}`, "/users-fisher", true)
	//testMatch(t, `/users/{userName:[\w\W]+`, "/users/fisher", true)
	testMatch(t, `/users/:userName`, "/users/fisher", false)
}

func TestRoute_MatchChildren(t *testing.T) {
	testMatchChildren(t, `/users`, "/users", true, "")
	testMatchChildren(t, `/users`, "/users/fisher", true, "/fisher")
	testMatchChildren(t, `/users`, "/users-fisher", true, "-fisher")
	testMatchChildren(t, `/users/{userName}`, "/users/fisher/profile", true, "/profile")
}

// Test match RouteMethod.
func testMatch(t *testing.T, pattern, path string, expected bool) {
	route := NewRoute(pattern, nil)
	matches := route.Match(path, *NewFakeRequest(), *NewFakeResponse())
	if matches != expected {
		t.Fatal("Failed to match path, expected [", expected, "]: ", pattern, path)
	}
}

func testMatchChildren(t *testing.T, pattern, path string, expected bool, expectedRemainingPath string) {
	route := NewRoute(pattern, nil)
	matched, remainingPath := route.MatchPrefixPath(path, *NewFakeRequest(), *NewFakeResponse())
	if matched != expected || expectedRemainingPath != remainingPath {
		t.Fatal("Failed to match children, expected [", expected, expectedRemainingPath, "]: ", pattern, path, remainingPath)
	}
}
