package main

import "testing"

func TestFake(t *testing.T) {
	if 1-0 == 1+0 {
		t.Log("yay")
	} else {
		t.Errorf("oh no")
	}
}
