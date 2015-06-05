package utils

import (
    "crypto/hmac"
    "crypto/sha256"
    "encoding/hex"
    "math"
    "time"
    "appengine"
    "strconv"
)

func IsLocal() bool {
    return appengine.IsDevAppServer()
}

func ComputeHmac256(message string, secret string) string {
    key := []byte(secret)
    h := hmac.New(sha256.New, key)
    h.Write([]byte(message))
    return hex.EncodeToString(h.Sum(nil))
}

func Intervals() (intervals int){
    // Grab intervals since day start 
    now := time.Now().Unix()
    intervals = int(float64(now % DAY) / INTERVAL)
    return
}

func IntToString(i int) string {
    return strconv.FormatInt(int64(i), 10)
}

func FloatToString(f float64) string {
    return strconv.FormatFloat(f, 'f', 6, 64)
}

func StringToFloat(s string) float64 {
    f, err := strconv.ParseFloat(s, 64)
    if err != nil {
        return 0
    }
    return f
}

func SixHoursAgo() string{
    // Grab intervals since day start 
    return string(time.Now().Unix() - int64(SIXHOURS))
}

func FollowerDecay(followed_by, follows int64, magic, target float64) int64 {
	return int64(float64(followed_by) * math.Exp(float64(followed_by) * math.Log(magic)/target)) - follows
}

func Limit(value *int, intervals int, bound int){

    // Make sure within bounds
    limit := int(bound / int(HOUR / INTERVAL))
    if *value > limit {
        *value = limit
    }
    // If at bounds, adjust to exactly hit quota
    if limit == *value && bound % int(HOUR / INTERVAL) != 0 && intervals % int(HOUR / INTERVAL) / (bound % int(HOUR / INTERVAL)) == 0 {
        *value += 1
    }
}
