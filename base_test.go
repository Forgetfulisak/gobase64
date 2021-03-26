package main

import (
	"testing"
)

func assertEqual(a, b string, t *testing.T) {
	if a != b {
		t.Fatal(a, "not equal to: ", b)
	}
}

func Test_encode(t *testing.T) {

	assertEqual(encodeString(""), "", t)
	assertEqual(encodeString("f"), "Zg==", t)
	assertEqual(encodeString("fo"), "Zm8=", t)
	assertEqual(encodeString("foo"), "Zm9v", t)
	assertEqual(encodeString("foob"), "Zm9vYg==", t)
	assertEqual(encodeString("fooba"), "Zm9vYmE=", t)
	assertEqual(encodeString("foobar"), "Zm9vYmFy", t)
}

func Test_Decode(t *testing.T) {

	assertEqual("", string(decodeString("")), t)
	assertEqual("f", string(decodeString("Zg==")), t)
	assertEqual("fo", string(decodeString("Zm8=")), t)
	assertEqual("foo", string(decodeString("Zm9v")), t)
	assertEqual("foob", string(decodeString("Zm9vYg==")), t)
	assertEqual("fooba", string(decodeString("Zm9vYmE=")), t)
	assertEqual("foobar", string(decodeString("Zm9vYmFy")), t)
}
