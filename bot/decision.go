package bot

func BasicDecision(s *Session, follows int, likes int, intervals int){
    // Round robin the hashtags. Allows for manual weighting eg: [#dog,#dog,#cute] 
    posts := GetPosts(s,s.GetHashtag(intervals))

    // Go from end to reduce collision
    i := 19
    for (likes > 0 || follows > 0) && i >= 0 {

        // Process likes
        if likes > 0 {
            go LikePosts(s, posts.Data[i].Id)
            likes--

        // Doing this seperately reaches larger audience
        // Never exceeds 12/11 at a given time
        }else if follows > 0 {
            go FollowUser(s, posts.Data[i].Id)
            follows--
        }

        // Decrement
        i--
    }
}

func IntelligentDecision(s *Session, follow int, likes int, intervals int) {
    //TODO: Smarter hashtags
//    posts := GetPosts(s,s.GetHashtag(intervals))
    //  gradient := s.
}


//func EvaluateProbability(post Post) float64 {
//    return Sigmoid(person,s.GetTheta())
//}