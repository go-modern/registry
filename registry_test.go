package registry_test

import (
	"runtime/debug"
	"testing"

	"github.com/go-modern/registry"
)

func TestRegistry(t *testing.T) {
	r := registry.New[string, string]("test")

	if r.Get("foo") != nil {
		t.Fatal("unexpected value")
	}

	r.Put("foo", "bar")

	if *r.Get("foo") != "bar" {
		t.Fatal("unexpected value")
	}

	if err := r.TryPut("foo", "bar"); err == nil {
		t.Fatal("expected error")
	}

	if err := r.TryPut("foo", "baz"); err == nil {
		t.Fatal("expected error")
	}

	if err := r.TryPut("bar", "baz"); err != nil {
		t.Fatal("unexpected error")
	}

	shouldPanic(t, func() { r.Put("foo", "baz") })

	shouldPanic(t, func() { r.Put("bar", "baz") })
}

func shouldPanic(t *testing.T, f func()) {
	defer func() {
		if r := recover(); r == nil {
			debug.PrintStack()
			t.Fatal("expected panic")
		}
	}()
	f()
}
