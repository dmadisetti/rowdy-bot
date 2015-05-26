package bot

import(
    "log"
    "bot/jobs"
    "bot/session"
    "bot/utils"
    "bot/http"
)

// Supervised learning -
//// Ideal person is a follower
//// Non ideal person is someone we follow, but doesn't follow back


// TODO: Clean this guy
func Learn(s *session.Session) string{

    if(s.SetLearning()){
        // New
        log.Println("Set up learning")
        status := http.GetStatus(s)
        s.SetLimits(int(status.Follows), int(status.Followed_by))
    }

    switch s.GetLearningStep() {

    case utils.APPRAISE:
        jobs.MinePeople(s)
        return StatusBar(s,"Mining Followers")

    case utils.SCORN:
        jobs.MinePeople(s)
        return StatusBar(s,"Mining Following")

    case utils.BUILD:
        // Logistic Regression
        // Get records and run
        go jobs.LogisticRegression(s)
        s.IncrementStep()
        return "* Running Logistic Regression"

    case utils.GOODTAGS:
        go jobs.CrawlTags(s, true)
        return StatusBar(s,"Finding Good Tags")

    case utils.BADTAGS:
        go jobs.CrawlTags(s, false)
        return StatusBar(s,"Finding Bad Tags")

    case utils.COMPUTETAGS:
        go jobs.WeightTags(s)
        return "* Ranking Tags"

    case utils.SHARE:
        go s.Share()
        //s.IncrementStep()
        return "Sharing"
    }
    
    return "Stop"
}

func StatusBar(s *session.Session, title string) (bar string) {
    bar = "    " + title + ":"
    BARSIZE := 100 - len(title)
    i := int(BARSIZE * s.GetState()/s.GetLimit())
    j := BARSIZE - i

    for i + j > 0 {
        if i > 0 {
            i--
            bar += "*"
        }else {
            j--
            bar += "-"
        }
    }
    bar += utils.IntToString(s.GetState()) + "/" + utils.IntToString(s.GetLimit())
    return
}
