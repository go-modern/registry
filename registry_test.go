package registry_test

import (
	"runtime/debug"
	"testing"

	"github.com/go-modern/registry"
)

func TestRegistry(t *testing.T) {
	r := registry.New[string, string]("test")

	if r.Load("foo") != nil {
		t.Fatal("unexpected value")
	}

	r.MustStore("foo", "bar")

	if *r.Load("foo") != "bar" {
		t.Fatal("unexpected value")
	}

	if err := r.ShouldStore("foo", "bar"); err == nil {
		t.Fatal("expected error")
	}

	if err := r.ShouldStore("foo", "baz"); err == nil {
		t.Fatal("expected error")
	}

	if err := r.ShouldStore("bar", "baz"); err != nil {
		t.Fatal("unexpected error")
	}

	shouldPanic(t, func() { r.MustStore("foo", "baz") })

	shouldPanic(t, func() { r.MustStore("bar", "baz") })
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
