package bot

import(
    "log"
    "strings"
    "net/url"
    "encoding/json"
)

// Could dynamically build these, but naw

// Get actions
func get(s *Session, url string, class interface{}) (*json.Decoder, interface{}) {
    response,err := s.Get(url)
    if err != nil {
        panic(err)
    }

   //Decode request
    decoder := json.NewDecoder(response.Body)
    err = decoder.Decode(&class)
    if err != nil {
        panic(err)
    }

    return decoder, class
}

func GetStatus(s *Session) (counts Counts){
    // var status Status
    _, face := get(s, "https://api.instagram.com/v1/users/" + s.GetId(), new(Status))
    if status, ok := face.(Status); ok {
        counts = status.Data.Counts
    } else {
        panic("GetStatus Broke")
    }
    return
}

func GetPosts(s *Session, hashtag string) (posts Posts){
    decoder, face := get(s, "https://api.instagram.com/v1/tags/" + hashtag +"/media/recent", new(Posts))
    if ps, ok := face.(Posts); ok {
        posts = ps
        posts.Next = getPagination(decoder)
    } else {
        panic("GetPosts Broke")
    }
    return posts
}

func GetUser(s *Session, id string) (user User){
    _, face := get(s, "https://api.instagram.com/v1/users/" + id, new(User))
    if u, ok := face.(User); ok {
        user = u
    } else {
        panic("GetPosts Broke")
    }
    return
}

func getPeople(s *Session, url string, id string) (users Users){
    decoder, face := get(s, url, new(Users))
    if us, ok := face.(Users); ok {
        users = us
        users.Next = getPagination(decoder)
    } else {
        panic("GetPosts Broke")
    }
    return users
}

func GetFollowing(s *Session, id string) Users{
    return getPeople(s, "https://api.instagram.com/v1/users/" + s.GetId() +"/follow", id)
}

func GetFollowers(s *Session, id string) Users{
    return getPeople(s, "https://api.instagram.com/v1/users/" + s.GetId() +"/followed-by", id)
}

func GetTag(s *Session, hashtag string) (tag Tag){
    _, face := get(s, "https://api.instagram.com/v1/tags/" + hashtag, new(Tag))
    if t, ok := face.(Tag); ok {
        tag = t
    } else {
        panic("GetPosts Broke")
    }
    return
}

// Post actions
func post(s *Session, url string, v url.Values) {
    response ,err := s.Post(url,v)
    if err != nil {
        panic(err)
    }
    log.Println(response)
}

func LikePosts(s *Session, id string) {
    v := url.Values{}
    post(s, "https://api.instagram.com/v1/media/"+id+"/likes", v)
}

func FollowUser(s *Session, id string){
    v := url.Values{}
    v.Set("action", "follow")
    post(s, "https://api.instagram.com/v1/users/"+ strings.Split(id,"_")[1] +"/relationship", v)
}

// Helper to grab next page
func getPagination(decoder *json.Decoder) string{
    var page Pagination
    err := decoder.Decode(&page)
    if err != nil {
        panic(err)
    }
    return page.Next_url
}