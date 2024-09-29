package pswdhash

import (
	"testing"
)

func TestHashSame(t *testing.T) {
	a := "ABC"

	aHash, err := HashPassword(a)
	if err != nil {
		t.Fatal("Could not hash password a")
	}
	if !VerifyPassword(a, aHash) {
		t.Fatal("Hash did not return the same (NOT GOOD!)")
	}
}

func TestHashDifference(t *testing.T) {
	a := "ABC"
	b := "CBA"

	aHash, err := HashPassword(a)
	if err != nil {
		t.Fatal("Could not hash password a")
	}

	bHash, err := HashPassword(b)
	if err != nil {
		t.Fatal("Could not hash password b")
	}

	if aHash == bHash {
		t.Fatal("Hash did return the same (NOT GOOD!)")
	}
}
