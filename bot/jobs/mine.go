package jobs

import(
	"bot/session"
	"bot/http"
	"bot/utils"
	"log"
)

func MinePeople(s *session.Session) {
    // Set up channel
    next := make(chan *http.Users)
    var batch http.Users

    if s.GetNext() == "" {
        if s.GetState() > 0 {
            s.IncrementStep()
        }
        if(s.GetLearningStep() == utils.APPRAISE){
            batch = http.GetFollowers(s)
        }else{
            batch = http.GetFollowing(s)
        }
    } else {
        batch = http.GetNext(s, s.GetNext())
    }

    go listen(s, next, 0, s.IncrementState())
    next <- &batch
}

func process(s *session.Session, users *http.Users, i int, follows float64){
    for i >= 0 {

        id := users.Data[i].Id
        user  := http.GetUser(s, id)

        log.Println(user)

        if(user.Data.Counts.Followed_by + user.Data.Counts.Follows > 0){
            //check follower records, if following and guy in other records, don't do anythin
            person := session.Person{
                Followers: float64(user.Data.Counts.Follows),
                Following: float64(user.Data.Counts.Followed_by),
                Posts: float64(user.Data.Counts.Media),
                Follows: !s.CheckCache(id),
            }

            // Because unset properties won't change, this should be fine
            if int(follows) == utils.SCORN {
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
func listen(s *session.Session, next chan *http.Users, calls int, follows float64) {
    for {
        select {
            case users := <-next:

                i := len(users.Data) - 1
                s.IncrementCount()
                go process(s, users, i, follows)

                close(next)
                if calls == utils.MAXPEOPLEGRAB {
                    s.SetNext(users.Pagination.Next_url)
                    return
                }

                var batch http.Users
                nxt := make(chan *http.Users)
                if users.Pagination.Next_url != "" {
                    log.Println("Getting another batch")
                    batch = http.GetNext(s, users.Pagination.Next_url)
                }else if follows == 0{ // follows == float64(s.GetLearningStep()) then have a array of functions
                    log.Println("Proceeding to next Step")
                    s.IncrementStep()
                    s.IncrementState()
                    batch = http.GetFollowing(s)
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
