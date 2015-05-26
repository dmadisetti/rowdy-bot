package bot

import(
    "net/http"
    "fmt"
    "strings"
    "html/template"
    "bot/session"
    "bot/utils"
    action "bot/http"
)

var t *template.Template

// Start er up!
func init(){
    NewHandler("/", mainHandle)
    NewHandler("/init", initHandle)
    NewHandler("/auth", authHandle)
    NewHandler("/process", processHandle)

    // For ML
    NewHandler("/learn", learningHandle)
    NewHandler("/update", updateHandle)
    NewHandler("/flush", flushHandle)
    NewHandler("/flushhashtag", flushHashtagHandle)

    // For testing
    NewHandler("/tag", tagHandle)
    NewHandler("/user", userHandle)
}

// Handles
func mainHandle(w http.ResponseWriter, r *http.Request, s *session.Session){
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

// Handle takes care of auth. Just for clean url
func initHandle(w http.ResponseWriter, r *http.Request, s *session.Session){
    http.Redirect(w,r,"/",302)
}

func authHandle(w http.ResponseWriter, r *http.Request, s *session.Session){
    s.SetHashtags(strings.Split(r.URL.Query()["hashtags"][0]," "))
    action.Authenticate(s,r.URL.Query()["code"][0])

    http.Redirect(w,r,"/",302)
}

func processHandle(w http.ResponseWriter, r *http.Request, s *session.Session){

    // Grab intervals since day start
    intervals := utils.Intervals()

    // Had some fancy math for peroidictiy. But
    // we could just brute force 100 per hour
    likes := int(utils.LIKES / utils.CALLS)
    utils.Limit(&likes, intervals, utils.LIKES)

    if !s.Usable() {
        fmt.Fprint(w, "Please set hashtags and authorize")
        return
    }

    // Follow ratio function where target is the desired
    // amount of followers.
    // e^(x*ln(magic)/target)
    // I wish could say there's some science behind why
    // we're doing this, but ultimately we just need a
    // decreasing function and some percentage of your
    // target feels right
    count := action.GetStatus(s)
    follows := int(utils.FollowerDecay(count.Followed_by,count.Follows,s.GetMagic(),s.GetTarget()))
    utils.Limit(&follows, intervals, utils.FOLLOWS)
    if follows < 0 {
        follows = 0
    }

    // Hang on channel otherwise jobs cut out
    done := make(chan bool)

    // Save status at midnight
    if intervals == 0 {
        go s.SetRecords(count.Followed_by,count.Follows)
    }
    if s.GetLearnt(){
        IntelligentDecision(s, follows, likes, intervals, done)
    } else {
        BasicDecision(s, follows, likes, intervals, done)
    }
    // Wait for finish. Defeats the purpose of aysnc, but only way to run in prod
    <-done
    fmt.Fprint(w, "Processing")
}

// Learning handle. Majority of logic in sentience.go
func learningHandle(w http.ResponseWriter, r *http.Request, s *session.Session){
    // If we could do this without being charged $10bijallion we could remove the conditional
    if utils.IsLocal() {
        fmt.Fprint(w, Learn(s))
        return
    }
    fmt.Fprint(w,"Must be local")
}

// Flush handle to kill all ML data
func flushHandle(w http.ResponseWriter, r *http.Request, s *session.Session){
    s.Flush()
    fmt.Fprint(w,  "Done Flushed")
}
// Flush handle to kill all ML data
func flushHashtagHandle(w http.ResponseWriter, r *http.Request, s *session.Session){
    go s.FlushEntity("Hashtag")
    fmt.Fprint(w,  "Done Flushed")
}

func updateHandle(w http.ResponseWriter, r *http.Request, s *session.Session){
    // Probs implement TOTP, potentially vulnerable to MTM
    if s.VerifiedUpdate(r.URL.Query()["hash"][0]){
        s.SetHashtags(strings.Split(r.URL.Query()["hashtags"][0]," "))
        s.ParseTheta(strings.Split(r.URL.Query()["theta"][0]," "))
        s.SetLearnt()
        fmt.Fprint(w,  "Updated")
    }else{
        fmt.Fprint(w,  "Not Verified")        
    }
}

// Just some testing endpoints
func tagHandle(w http.ResponseWriter, r *http.Request, s *session.Session){
    tag := action.GetTag(s, r.URL.Query()["hashtag"][0])
    fmt.Fprint(w, tag.Data.Media_count)
}

// Snoop Doggy Dog
// http://127.0.0.1:8080/user?user=1574083
func userHandle(w http.ResponseWriter, r *http.Request, s *session.Session){
    user := action.GetUser(s, r.URL.Query()["user"][0])
    fmt.Fprint(w, user.Data.Counts.Followed_by)
}
