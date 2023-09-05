package main

import (
	"crypto/hmac"
	"crypto/sha256"
)

func HashPassword(password string, pass_secret []byte) []byte {
	h := hmac.New(sha256.New, pass_secret)
	h.Write([]byte(password))
	return h.Sum(nil)
}
