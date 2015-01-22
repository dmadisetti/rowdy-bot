package bot

import (
    "crypto/hmac"
    "crypto/sha256"
    "encoding/hex"
    "math"
)

func ComputeHmac256(message string, secret string) string {
    key := []byte(secret)
    h := hmac.New(sha256.New, key)
    h.Write([]byte(message))
    return hex.EncodeToString(h.Sum(nil))
}

func FollowerDecay(count Counts, magic float64, target float64) int64 {
	return int64(float64(count.Followed_by) * math.Exp(float64(count.Followed_by) * math.Log(magic)/target))
}