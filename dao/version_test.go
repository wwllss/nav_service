package dao

import "testing"

func TestVersion_After(t *testing.T) {
	v := Version{
		Version: "4.88.0",
	}
	if v.Before("4.87.0") {
		t.Fatal("error")
	}
}
