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

	if err := r.Store("foo", "bar"); err == nil {
		t.Fatal("expected error")
	}

	if err := r.Store("foo", "baz"); err == nil {
		t.Fatal("expected error")
	}

	if err := r.Store("bar", "baz"); err != nil {
		t.Fatal("unexpected error")
	}

	shouldPanic(t, func() { r.MustStore("foo", "baz") })

	shouldPanic(t, func() { r.MustStore("bar", "baz") })

	v := r.Default()
	if v != nil {
		t.Fatal("expected nil")
	}

	v, err := r.Init("foo")
	if err != nil {
		t.Fatal("unexpected error")
	}
	if *v != "bar" {
		t.Fatal("unexpected value")
	}

	v = r.Default()
	if *v != "bar" {
		t.Fatal("unexpected value")
	}
	_, err = r.Init("foo")
	if err == nil {
		t.Fatal("expected error")
	}
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
