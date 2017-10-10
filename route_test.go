package wstf

import (
	"testing"
)

func TestRoute_Match(t *testing.T) {
	testMatch(t, `/users`, "/users", true)
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
	matches := route.Match(path, Request{}, Response{})
	if matches != expected {
		t.Fatal("Failed to match, expected [", expected, "]: ", pattern, path)
	}
}
func testMatchChildren(t *testing.T, pattern, path string, expected bool, expectedRemainingPath string) {
	route := NewRoute(pattern, nil)
	matched, remainingPath := route.MatchPrefixPath(path, Request{}, Response{})
	if matched != expected || expectedRemainingPath != remainingPath {
		t.Fatal("Failed to match, expected [", expected, expectedRemainingPath, "]: ", pattern, path, remainingPath)
	}
}
