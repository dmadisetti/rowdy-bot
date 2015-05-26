package bot

import(
    "bot/session"
    "bot/utils"
    "bot/http"
    "strings"
)

func BasicDecision(s *session.Session, follows int, likes int, intervals int, done chan bool){
    // Round robin the hashtags. Allows for manual weighting eg: [#dog,#dog,#cute] 
    posts := http.GetPosts(s,s.GetHashtag(intervals))

    // Go from end to reduce collision
    i := 19
    for (likes > 0 || follows > 0) && i >= 0 {

        // Process likes
        if likes > 0 {
            go http.LikePosts(s, posts.Data[i].Id)
            likes--

        // Doing this seperately reaches larger audience
        // Never exceeds 12/11 at a given time
        }else if follows > 0 {
            go http.FollowUser(s, posts.Data[i].Id)
            follows--
        }

        // Decrement
        i--
    }

    // Indicate doneness
    done <- true
}

func IntelligentDecision(s *session.Session, follows int, likes int, intervals int,  done chan bool) {

    // Still do round robin, but this time the hashtags are smart
    posts := http.GetPosts(s,s.GetHashtag(intervals))
    next := make(chan *http.Posts)
    grp := make(chan *group)
    count := 0
    calls := 0
    go sort(s, grp, follows, likes, &calls, &count, done)
    go listen(s, grp, next, &calls, &count)
    next <- &posts
}

// Async heapsort, hope it works
func sort(s *session.Session, next chan *group, follows, likes int, calls, total *int, done chan bool) {
    var instances []group
    count := 0
    x := 0
    min := 1.1
    for {
        select {
            case instance := <-next:

                x++
                // Catches up and thus done
                if x == *total && *calls == utils.MAXPOSTGRAB {
                    i := 0
                    for (likes > 0 || follows > 0){

                        // Highest value for follows then do likes
                        if follows > 0 {                            
                            //http.FollowUser(s, instances[i].id)
                            follows--
                        }else if likes > 0 {
                            //http.LikePosts(s, instances[i].id)
                            likes--
                        }
                        i++
                    }
                    s.FlushCache()
                    done <- true
                    close(next)
                    return
                }

                if instance.id == "continue" || (instance.value <= min && count == follows + likes) {
                    continue
                }

                if min < instance.value {
                    if count == follows + likes {
                        min = instance.value
                    }
                } else {
                    if count < follows + likes {
                        min = instance.value
                    }                    
                }

                if count < follows + likes {
                    instances = append(instances, *instance)
                    count += 1
                } else {
                    instances[count - 1] = *instance
                }

                // Bubble sort
                for i := count - 2; i >= 0; i-- {
                    if instance.value > instances[i].value {
                        holder := instances[i]
                        instances[i] = *instance
                        instances[i + 1] = holder
                    } else {
                        break
                    }
                }
        }
    }
}

type group struct {
    value float64
    id string
    user string
}

// Async set up multi calls
func listen(s *session.Session, grp chan *group, next chan *http.Posts, calls, count *int) {
    for {
        select {
            case posts := <-next:

                i := len(posts.Data) - 1
                *count += len(posts.Data)
                go process(s, posts, i, grp)

                close(next)
                if *calls == utils.MAXPOSTGRAB || posts.Pagination.Next_url == "" {
                    return
                }

                var batch http.Posts
                nxt := make(chan *http.Posts)
                batch = http.GetNextPost(s, posts.Pagination.Next_url)

                *calls += 1
                go listen(s, grp, nxt, calls, count)
                nxt <- &batch
                return
        }
    }
}

func process(s *session.Session, posts *http.Posts, i int, grp chan *group){
    for i >= 0 {

        id := strings.Split(posts.Data[i].Id,"_")[1]
        if posts.Data[i].User_has_liked || s.CheckCache(id) || http.IsFollowing(s,id){
            // Try to add to channel and stop if done
            grp <- &group{
                id:"continue",
            }
            i--
            continue
        }
        user  := http.GetUser(s, id)
        // Create perosn to get value
        person := session.Person{
            Followers: float64(user.Data.Counts.Follows),
            Following: float64(user.Data.Counts.Followed_by),
            Posts: float64(user.Data.Counts.Media),
        }

        grp <- &group{
            id:id,
            value: person.Sigmoid(s.GetTheta()),
            user: posts.Data[i].User.Username,
        }

        i--
    }
}
