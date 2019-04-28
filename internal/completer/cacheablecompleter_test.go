package completer

import (
	"testing"
)

func TestQueryToKey(t *testing.T) {
	q := &Query{
		Term:   "Moscow",
		Locale: "En",
		Types:  []string{"zaw", "dfg", "asd", "zxc"},
	}

	cacheKey := queryToCacheKey(q)
	expected := "asddfgzawzxcMoscowEn"
	if actual := string(cacheKey); actual != expected {
		t.Fatalf("cache key is invalid. Expected: %s, actual: %s", expected, actual)
	}
}
