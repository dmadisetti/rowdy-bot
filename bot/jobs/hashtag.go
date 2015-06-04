package jobs

import(
    "log"
    "bot/session"
    "bot/utils"
    "bot/http"
)

func CrawlTags(s *session.Session, positive bool){
    keys := s.GetPeopleCursor(positive, utils.MAXREQUESTS * s.GetState())
    s.IncrementState()
    z := 0
    size := 0
    for {
        key, err := keys.Next(nil)
        z += 1
        if err != nil {
            if z < utils.MAXREQUESTS {
                s.IncrementStep()
            }
            break // No further entities match the query.
        }
        media := http.GetMedia(s,key.StringID()).Data
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
        }
    }
    s.IncrementSize(size,positive)
    s.StopProcessing()
}

func processTags(s *session.Session, weight float64, tags []string){
    weight /= float64(len(tags))
    for _, next := range tags {
        h := s.Hashtag(next)
        h.Value += weight
        s.SaveHashtag(h)
    }
}

func WeightTags(s *session.Session){
    
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
    s.StopProcessing()
}