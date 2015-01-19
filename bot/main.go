package bot

import(
    "net/http"
    "fmt"
    "time"
    "math"
    "strings"
    "log"
    "html/template"
)

var t *template.Template
const DAY int64 = 60000 * 60 * 24
const INTERVAL float64 = 60000 * 5

// Good enough
const SQRT3OVER2 float64 = 0.86602540378 // math.Sqrt(3)/2

// Start er up!
func init(){
    NewHandler("/", mainHandle)
    NewHandler("/auth", authHandle)
    NewHandler("/hashtag", hashtagHandle)
    NewHandler("/process", processHandle)
}

// Handles
func mainHandle(w http.ResponseWriter, r *http.Request, s *Session){
    t, e := template.ParseGlob("templates/the.html")
    if e != nil {
        fmt.Fprint(w, e)
        return
    }
    // render with records
    err := t.Execute(w, s)
    if err !=nil{
        panic(err)
    }
}

func authHandle(w http.ResponseWriter, r *http.Request, s *Session){
    s.SetAuth(r.URL.Query()["code"][0])
    http.Redirect(w,r,"/",302)
}

func hashtagHandle(w http.ResponseWriter, r *http.Request, s *Session) {
    s.SetHashtags(strings.Split(r.URL.Query()["hashtags"][0],","))
    http.Redirect(w,r,"/",302)
}

func processHandle(w http.ResponseWriter, r *http.Request, s *Session){

    // Grab intervals since day start 
    now := time.Now().Unix()
    intervals := math.Floor(float64(now % DAY) / INTERVAL)

    // Our golden function. 
    // (cos((pi*x/144)-42) + sqrt(3)/2)/(1+sqrt(3)/2)
    // Cyclical cos function adjusted to represent the 
    // day in 288 parts (intervals of 5 minutes) and for
    // The lowest part of the day to be around 1AM-5AM EST
    // While highest is 4pm with a theoretical peak 
    // of 8.346 posts to fit under the limit of 100post/hr 
    // (You can do the riemann sum yourself)
    likes := int(((math.Cos((math.Pi*intervals/144)-42) + SQRT3OVER2)/(1 + SQRT3OVER2)) * 8.346)


    // Round robin the hashtags. Allows for manual weighting eg: [#dog,#dog,#cute]
    if len(s.Settings.Hashtags) == 0 {
        fmt.Fprint(w, "Please set hashtags")
        return
    }
    posts := GetPosts(s,s.Settings.GetHashtag(intervals))

    // Follow ratio function where target is the desired
    // amount of followers. 
    // e^(x*ln(magic)/target)
    // I wish could say there's some science behind why 
    // we're doing this, but ultimately we just need a
    // decreasing function and some percentage of your 
    // target feels right
    count := GetFollowing(s)
    follows := int64(float64(count.Followed_by) * math.Exp(float64(-count.Followed_by) * math.Log(s.Settings.Magic)/s.Settings.Target))
    follows -= count.Follows // New follows

    // Save status at midnight
    if intervals == 0 {
        s.SetRecords(count)
    }

    // Go from end to reduce collision
    i := 19
    for likes > 0 || follows > 0 {

        if likes > 0 {
            LikePosts(s, posts[i].Id)
        }
        if follows > 0 {
            FollowUser(s, posts[i].Id)
        }

        log.Println(posts[i])

        // Decrement
        likes--
        follows--
        i--
    }
}
