package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/imroc/req"
)

func authHeaders(body, method, path string) req.Header {
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	return req.Header{
		"Content-Type":         "Application/JSON",
		"CB-ACCESS-KEY":        os.Getenv("API_KEY"),
		"CB-ACCESS-SIGN":       sign(body, strings.ToUpper(method), path, timestamp),
		"CB-ACCESS-TIMESTAMP":  timestamp,
		"CB-ACCESS-PASSPHRASE": os.Getenv("API_PASSPHRASE"),
	}
}

func sign(body, method, path, timestamp string) string {
	secret, _ := base64.StdEncoding.DecodeString(os.Getenv("API_SIGNATURE"))
	envelop := hmac.New(sha256.New, []byte(secret))
	message := fmt.Sprintf("%s%s%s%s", timestamp, method, path, body)
	envelop.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(envelop.Sum(nil))
}
