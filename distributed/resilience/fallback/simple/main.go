package main

import (
	"errors"
	"fmt"
)

var staleCache = map[string]string{"user:1": "alice (cached)"}

func fetchUser(id string) (string, error) {
	return "", errors.New("upstream unavailable")
}

func getUser(id string) string {
	if v, err := fetchUser(id); err == nil {
		return v
	}
	if v, ok := staleCache[id]; ok {
		fmt.Printf("[fallback] returning stale cache for %s\n", id)
		return v
	}
	fmt.Printf("[fallback] using default for %s\n", id)
	return "anonymous"
}

func main() {
	fmt.Println(getUser("user:1"))  // cache hit
	fmt.Println(getUser("user:99")) // default
}
