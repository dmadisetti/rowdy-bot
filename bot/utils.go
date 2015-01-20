package bot

import (
    "bytes"
    "crypto/hmac"
    "crypto/sha256"
    "encoding/hex"
)

func Buffer(s string) *Buffer{
	return bytes.NewBufferString(s)
}

func ComputeHmac256(message string, secret string) string {
    key := []byte(secret)
    h := hmac.New(sha256.New, key)
    h.Write([]byte(message))
    return hex.EncodeToString(h.Sum(nil))
}