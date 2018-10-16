// + build functional
package func_test

import "testing"

func TestFooEqualsBar(t *testing.T) {
	got := "foobar"
	expected := "foobar"
	if got != expected {
		t.Fatalf("got %s, expected %s", got, expected)
	}
}

func TestFooDoesNotEqualBar(t *testing.T) {
	got := "foobar"
	expected := "foobar"
	if got == expected {
		t.Fatalf("got %s, expected not %s", got, expected)
	}
}
