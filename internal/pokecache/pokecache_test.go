package pokecache

import (
	"testing"
	"time"
)

func TestAddGet(t *testing.T) {
	cache := NewCache(time.Hour)
	cache.Add("first", []byte("original"))
	cache.Add("second", []byte("other"))
	cache.Add("first", []byte("updated"))

	for _, tc := range []struct {
		key   string
		want  string
		found bool
	}{
		{key: "first", want: "updated", found: true},
		{key: "second", want: "other", found: true},
		{key: "missing", found: false},
	} {
		t.Run(tc.key, func(t *testing.T) {
			got, ok := cache.Get(tc.key)
			if ok != tc.found || string(got) != tc.want {
				t.Fatalf("Get(%q) = %q, %v; want %q, %v", tc.key, got, ok, tc.want, tc.found)
			}
			if !ok && got != nil {
				t.Fatal("missing key should return nil")
			}
		})
	}
}

func TestReap(t *testing.T) {
	for _, tc := range []struct {
		name  string
		age   time.Duration
		found bool
	}{
		{name: "fresh", age: time.Second, found: true},
		{name: "expired", age: 10 * time.Second, found: false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cache := NewCache(time.Hour)
			cache.Add("page", []byte("locations"))
			cache.reap(time.Now().Add(tc.age), 5*time.Second)
			_, ok := cache.Get("page")
			if ok != tc.found {
				t.Fatalf("Get after reaping = %v; want %v", ok, tc.found)
			}
		})
	}
}

func TestReapLoop(t *testing.T) {
	cache := NewCache(20 * time.Millisecond)
	cache.Add("page", []byte("locations"))
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		if _, ok := cache.Get("page"); !ok {
			return
		}
		time.Sleep(5 * time.Millisecond)
	}
	t.Fatal("expected the background loop to remove the expired entry")
}
