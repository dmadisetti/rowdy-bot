package bot

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

func FollowerDecay(count Counts, magic float64, target float64) int64 {
	return int64(float64(count.Followed_by) * math.Exp(float64(count.Followed_by) * math.Log(magic)/target)) - count.Follows
}

func Limit(value *int, intervals int, bound int){

    // Make sure within bounds
    limit := int(bound / int(HOUR / INTERVAL))
    if *value > limit {
      *value = limit
    }

    // If at bounds, adjust to exactly hit quota
    if limit == *value && intervals % int(HOUR / INTERVAL) / (bound % int(HOUR / INTERVAL)) == 0 {
        *value += 1
    }
}

func Y(person Person) (y float64) {
    y = 1.0
    if person.Follows {
        y = 0.0
    }
    return
}

func Sigmoid(person Person, gradient []float64) float64{
	f := gradient[0] +  person.Followers * gradient[1] + person.Following * gradient[2] + person.Posts * gradient[3]
	return 1.0/(1.0 + math.Exp(-f))
}

func J(person Person, gradient []float64) float64{
    y := Y(person)
	h := Sigmoid(person, gradient)
	q := y * math.Log(h) + (1 - y) * math.Log(1 - h)
	return q
}
