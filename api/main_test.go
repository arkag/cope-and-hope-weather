package main

import "testing"

func TestPlaceholder(t *testing.T) {
	got := true
	want := true
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}
