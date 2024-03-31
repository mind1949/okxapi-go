package okxapigo

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"strconv"
)

// Sign calc signature
func Sign(timestamp int64, api Api) string {
	const method = http.MethodGet
	const requestPath = "/users/self/verify"

	hash := hmac.New(sha256.New, []byte(api.Secretkey))
	for _, s := range []string{
		strconv.FormatInt(timestamp, 10),
		method,
		requestPath,
	} {
		hash.Write([]byte(s))
	}
	sum := hash.Sum(nil)
	sign := base64.StdEncoding.EncodeToString(sum)
	return sign
}
