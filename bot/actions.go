package bot

import(
    "log"
    "strings"
    "net/url"
    "encoding/json"
)

func GetFollowing(s *Session) (count Counts){

    response,err := s.Get("https://api.instagram.com/v1/users/" + s.Settings.GetId())
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

    log.Println(status)

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

func LikePosts(s *Session, id string){
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

// func CommentPost(r *http.Request, id string){
//     c := appengine.NewContext(r)
//     client := urlfetch.Client(c)
// 
//     v := url.Values{}
//     v.Set("access_token", "Token")
//     v.Add("text", "woof!")
// 
//     response,err := client.PostForm("https://api.instagram.com/v1/media/"+id+"/comments",v)
//     if err != nil {
//         panic(err)
//     }
// 
//     log.Println(response)
// }
