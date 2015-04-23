package bot

import(
    "log"
    "strings"
    "net/url"
    "encoding/json"
)

// Could dynamically build these, but naw

// Get actions
func GetFollowing(s *Session) (count Counts){

    response,err := s.Get("https://api.instagram.com/v1/users/" + s.GetId())
    if err != nil {
        panic(err)
    }

   //Decode request
    var status Status
    decoder := json.NewDecoder(response.Body)
    err = decoder.Decode(&status)
    if err != nil {
        panic(err)
    }

    count = status.Data.Counts
    return
}

func GetPosts(s *Session, hashtag string) []Post{

    response,err := s.Get("https://api.instagram.com/v1/tags/" + hashtag +"/media/recent")
    if err != nil {
        panic(err)
    }

    log.Println(response)

    //Decode request
    var posts Posts
    decoder := json.NewDecoder(response.Body)
    err = decoder.Decode(&posts)
    if err != nil {
        panic(err)
    }

    return posts.Data
}

func GetUser(s *Session, id string) User{

    response,err := s.Get("https://api.instagram.com/v1/users/" + id)
    if err != nil {
        panic(err)
    }

    //Decode request
    var user User
    decoder := json.NewDecoder(response.Body)
    err = decoder.Decode(&user)
    if err != nil {
        panic(err)
    }

    return user

}


func GetTag(s *Session, hashtag string) Tag{

    response,err := s.Get("https://api.instagram.com/v1/tags/" + hashtag)
    if err != nil {
        panic(err)
    }

    //Decode request
    var tag Tag
    decoder := json.NewDecoder(response.Body)
    err = decoder.Decode(&tag)
    if err != nil {
        panic(err)
    }

    return tag

}

// Post actions
func LikePosts(s *Session, id string) {
    v := url.Values{}

    response ,err := s.Post("https://api.instagram.com/v1/media/"+id+"/likes",v)
    if err != nil {
        panic(err)
    }
    log.Println(response)
}

func FollowUser(s *Session, id string){
    v := url.Values{}
    v.Set("action", "follow")

    response,err := s.Post("https://api.instagram.com/v1/users/"+ strings.Split(id,"_")[1] +"/relationship",v)
    if err != nil {
        panic(err)
    }
    log.Println(response)
}