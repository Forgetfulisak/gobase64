package main

import (
	"log"
	"testing"
)

func assertEqual(a, b string) {
	if a != b {
		log.Fatalln(a, "not equal to: ", b, len(a), len(b))
	}
}

func Test_encode(t *testing.T) {

	assertEqual(encodeString(""), "")
	assertEqual(encodeString("f"), "Zg==")
	assertEqual(encodeString("fo"), "Zm8=")
	assertEqual(encodeString("foo"), "Zm9v")
	assertEqual(encodeString("foob"), "Zm9vYg==")
	assertEqual(encodeString("fooba"), "Zm9vYmE=")
	assertEqual(encodeString("foobar"), "Zm9vYmFy")
}

func Test_Decode(t *testing.T) {

	assertEqual("", string(decodeString("")))
	assertEqual("f", string(decodeString("Zg==")))
	assertEqual("fo", string(decodeString("Zm8=")))
	assertEqual("foo", string(decodeString("Zm9v")))
	assertEqual("foob", string(decodeString("Zm9vYg==")))
	assertEqual("fooba", string(decodeString("Zm9vYmE=")))
	assertEqual("foobar", string(decodeString("Zm9vYmFy")))
}
