package bot

import(
    "log"
    "math"
)

// Supervised learning -
//// Ideal person is a follower
//// Non ideal person is someone we follow, but doesn't follow back

func process(s *Session, users *Users, i int, follows float64){
    for i >= 0 {

        id := users.Data[i].Id
        user  := GetUser(s, id)

        log.Println(user)

        if(user.Data.Counts.Followed_by + user.Data.Counts.Follows > 0){
            //check follower records, if following and guy in other records, don't do anythin
            person := Person{
                Followers: float64(user.Data.Counts.Follows),
                Following: float64(user.Data.Counts.Followed_by),
                Posts: float64(user.Data.Counts.Media),
                Follows: !s.CheckCache(id),
            }

            // Because unset properties won't change, this should be fine
            if int(follows) == SCORN {
                person.Followed = true
                person.Follows = !person.Follows
            }

            // Add to variable and to Keys 
            s.PutPerson(person, id)
        }

        // Decrement
        i--
    }

    // Catches up and thus done
    if s.FinishedCount() {
        s.SavePeople()
    }
}

// Async set up multi calls
func listen(s *Session, next chan *Users, calls int, follows float64) {
    log.Println(MAXPEOPLEGRAB)
    for {
        select {
            case users := <-next:

                i := len(users.Data) - 1
                s.IncrementCount()
                go process(s, users, i, follows)

                close(next)
                if calls == MAXPEOPLEGRAB {
                    s.SetNext(users.Pagination.Next_url)
                    return
                }

                var batch Users
                nxt := make(chan *Users)
                if users.Pagination.Next_url != "" {
                    log.Println("Getting another batch")
                    batch = GetNext(s, users.Pagination.Next_url)
                }else if follows == 0{ // follows == float64(s.GetLearningStep()) then have a array of functions
                    log.Println("Proceeding to next Step")
                    s.IncrementStep()
                    s.IncrementState()
                    batch = GetFollowing(s)
                    follows = float64(s.GetLearningStep())
                } else {
                    s.SetNext("")
                    return
                }

                go listen(s, nxt, calls + 1, follows)
                nxt <- &batch
                return
        }
    }
}

// TODO: Clean this guy
func Learn(s *Session) string{

    if(s.SetLearning()){
        // New
        log.Println("Set up learning")
        status := GetStatus(s)
        s.SetLimits(int(status.Follows), int(status.Followed_by))
    }

    switch s.GetLearningStep() {

    case APPRAISE:
        minePeople(s)
        return StatusBar(s,"Mining Followers")

    case SCORN:
        minePeople(s)
        return StatusBar(s,"Mining Following")

    case BUILD:
        // Logistic Regression
        // Get records and run
        go logisticRegression(s)
        s.IncrementStep()
        return "* Running Logistic Regression"

    case GOODTAGS:
        go CrawlTags(s, true)
        return StatusBar(s,"Finding Good Tags")

    case BADTAGS:
        go CrawlTags(s, false)
        return StatusBar(s,"Finding Bad Tags")

    case COMPUTETAGS:
        go WeightTags(s)
        return "* Ranking Tags"

    case SHARE:
        go s.Share()
        // s.IncrementStep()
        return "Stop" //"Sharing"
    }
    
    return "Stop"
}

func StatusBar(s *Session, title string) (bar string) {
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
    bar += IntToString(s.GetState()) + "/" + IntToString(s.GetLimit())
    return
}


func minePeople(s *Session) {
    // Set up channel
    next := make(chan *Users)
    var batch Users

    if s.GetNext() == "" {
        if s.GetState() > 0 {
            s.IncrementStep()
        }
        if(s.GetLearningStep() == APPRAISE){
            batch = GetFollowers(s)
        }else{
            batch = GetFollowing(s)
        }
    } else {
        batch = GetNext(s, s.GetNext())
    }

    go listen(s, next, 0, s.IncrementState())
    next <- &batch
}

func CrawlTags(s *Session, positive bool){
    keys := s.GetPeopleCursor(positive, MAXREQUESTS * s.GetState())
    s.IncrementState()
    z := 0
    size := 0
    for {
        key, err := keys.Next(nil)
        z += 1
        if err != nil {
            if z < MAXREQUESTS {
                s.IncrementStep()
            }
            break // No further entities match the query.
        }
        media := GetMedia(s,key.StringID()).Data
        captured := len(media) - 1
        for i := 0; i < 3 && i < captured; i++ {
            tagCount := len(media[i].Tags)

            lim := 5
            for j, tag := range media[i].Tags {
                if tag == "" {
                    continue
                }
                if j >= lim{
                    break
                }
                h := s.Hashtag(tag)
                log.Println(tag)
                for k := 0; k < lim && k < tagCount; k++ {
                    if j == k {
                        continue
                    }
                    if tag == media[i].Tags[k] {
                        lim += 1
                        media[i].Tags[k] = ""
                        continue
                    }
                    if positive{
                        h.Beneficiaries = append(h.Beneficiaries, media[i].Tags[k])
                    }else {
                        h.Victims = append(h.Beneficiaries, media[i].Tags[k])
                    }
                }
                s.SaveHashtag(h)
                size += 1
            }
            log.Println("------------------------------")
        }
    }
    s.IncrementSize(size,positive)
}

func processTags(s *Session, weight float64, tags []string){
    weight /= float64(len(tags))
    for _, next := range tags {
        h := s.Hashtag(next)
        h.Value += weight
        s.SaveHashtag(h)
    }
}

func WeightTags(s *Session){
    
    // Stop multiprocessing
    if s.GetState() > 0 {
        return
    }
    s.IncrementState()

    tags := s.GetHashtagCursor()
    goodWeight := 1.0/s.GetHashtagSize(true)
    badWeight  := -1.0/s.GetHashtagSize(false)

    for {
        key, err := tags.Next(nil)
        if err != nil {
            log.Println(err)
            break // No further entities match the query.
        }
        hashtag := s.Hashtag(key.StringID())

        log.Println(hashtag.Name)

        processTags(s, goodWeight, hashtag.Beneficiaries)
        processTags(s, badWeight,  hashtag.Victims)
    }

    s.SetTopTags()

    // Move on
    s.IncrementStep()
}


func logisticRegression(s *Session){

    // Grab all people because of many iterations
    people := s.GetPeople()
    objective := Objective{
        People: people,
        Lambda: s.GetLambda(),
        Alpha: s.GetAlpha(),
        Size: float64(len(people)),
    }

    start := []float64{1, 1, 1, 0} // Bias, Following, Followers, Posts
    minimum := Minimize(objective, start)
    log.Println(minimum)

    s.SetTheta(minimum)
}

func Minimize(objective Objective, thetas []float64) []float64{

    nthetas := objective.EvaluateGradient(thetas)
    next := objective.EvaluateFunction(nthetas)
    value := objective.EvaluateFunction(thetas)

    i := 0

    for (math.IsNaN(value) || value >= next) && i < 1000 {
        // log.Println(value)
        thetas = nthetas
        value = next
        nthetas = objective.EvaluateGradient(thetas)
        next = objective.EvaluateFunction(nthetas)
        i += 1
    }

    return thetas
}

type Hashtag struct{
    Name string
    Value float64
    Beneficiaries []string
    Victims []string
}

type Person struct{
    Following float64
    Followers float64
    Posts float64
    Follows bool // 1 or 0
    Followed bool
}

type Objective struct{
    People []Person
    Lambda float64
    Alpha float64
    Size float64
}

// Cost Function
func (o Objective) EvaluateFunction(thetas []float64) float64{

    sum := 0.0
    for _, person := range o.People {
        sum += J(person,thetas)
    }
    //sum += -(o.Lambda/2) * thetas[3] * thetas[3]
    return -sum/o.Size
}

// Cost Function derivative
func (o Objective) EvaluateGradient(thetas []float64) (gradient []float64) {
    gradient = make([]float64, len(thetas))

    gradient[0] = o.CostD(0, thetas)
    gradient[1] = o.CostD(1, thetas)
    gradient[2] = o.CostD(2, thetas)
    //gradient[3] = o.CostD(3, thetas)

    return
}

func (o Objective) CostD(i int, thetas []float64) float64{

    sum := 0.0

    switch i {
    case 0:
        for _, person := range o.People {
            sum += (Sigmoid(person, thetas) - Y(person))
        }
        break;
    case 1:
        for _, person := range o.People {
            sum += (Sigmoid(person, thetas) - Y(person)) * person.Followers
        }
        break;
    case 2:
        for _, person := range o.People {
            sum += (Sigmoid(person, thetas) - Y(person)) * person.Following
        }
        break;
    case 3:
        for _, person := range o.People {
            sum += (Sigmoid(person, thetas) - Y(person)) * person.Posts
        }
        sum += o.Lambda * thetas[3]
        break;
    default:
        break;
    }
    return thetas[i] - o.Alpha * sum
}
