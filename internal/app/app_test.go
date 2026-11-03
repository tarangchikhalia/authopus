package app

import "testing"

func TestGreeting(t *testing.T) {
	if got := Greeting(); got != "Hello World!" {
		t.Fatalf("Greeting() = %q, want %q", got, "Hello World!")
	}
}
