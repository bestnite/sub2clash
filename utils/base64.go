package utils

import (
	"encoding/base64"
	"strings"
)

func DecodeBase64(s string, urlSafe bool) (string, error) {
	s = strings.TrimSpace(s)
	if len(s)%4 != 0 {
		s += strings.Repeat("=", 4-len(s)%4)
	}
	var decodeStr []byte
	var err error
	if urlSafe {
		decodeStr, err = base64.URLEncoding.DecodeString(s)
	} else {
		decodeStr, err = base64.StdEncoding.DecodeString(s)
	}
	if err != nil {
		return "", err
	}
	return string(decodeStr), nil
}

func EncodeBase64(s string, urlSafe bool) string {
	if urlSafe {
		return base64.URLEncoding.EncodeToString([]byte(s))
	}
	return base64.StdEncoding.EncodeToString([]byte(s))
}
