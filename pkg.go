package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
)

func getEnv[T typ](key string, def T) T {
	var res interface{} = def

	switch reflect.TypeOf(def).Kind() {
	case reflect.String:
		if !isEmpty(os.Getenv(key)) {
			res = os.Getenv(key)
		}

	case reflect.Uint32:
		if !isEmpty(os.Getenv(key)) {
			rs, err := strconv.Atoi(os.Getenv(key))
			if err != nil {
				log.Println("key env:", key, "error convert value, to set default value")
			} else {
				res = uint32(rs)
			}
		}

	}

	return res.(T)
}

func isEmpty(val string) bool {
	return strings.TrimSpace(val) == ""
}

// Sign is using the data and the secret to compute a HMAC(SHA256) to sign the body of the request.
// so the webhook can use this signature to verify that no data have been compromised.
func FFSignature(payloadBody []byte, secretToken []byte) string {
	mac := hmac.New(sha256.New, secretToken)
	_, _ = mac.Write(payloadBody)
	expectedMAC := mac.Sum(nil)
	return "sha256=" + hex.EncodeToString(expectedMAC)
}

// GetSignature returns an HMAC-SHA256 signature encoded in base64
func GetSignature(secret []byte, messages []byte) string {
	h := hmac.New(sha256.New, secret)
	h.Write(messages)

	return hex.EncodeToString(h.Sum(nil))
}

// CompareSignatures is used to compare given and reference signature using time constant algorithm
func CompareSignatures(given string, reference string) (bool, error) {
	h1, err := hex.DecodeString(given)
	if err != nil {
		return false, err
	}

	h2, err := hex.DecodeString(reference)
	if err != nil {
		return false, err
	}

	return subtle.ConstantTimeCompare(h1, h2) == 1, nil
}
