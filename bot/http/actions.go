package http

import(
    "strings"
    "io/ioutil"
    "net/url"
    "encoding/json"
    "bot/session"
    "bot/utils"
)

// Could dynamically build these, but naw- generated code ever feels nice IM(Humble)O
// I'm sure there's a nicer interface implementation, 
// but copy and pasting was all too easy

// Authenticate
func Authenticate(s *session.Session, code string){
    decoder := s.Auth(code)
 
    //Decode request
    var auth Auth
    err := decoder.Decode(&auth)
    if err != nil {
        panic(err)
    }

    s.SetAuth(auth.Access_token, auth.User.Id)

}


// Get actions
func GetStatus(s *session.Session) (count Counts){

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

func GetMedia(s *session.Session, id string) Posts{
    params := map[string]string{"MIN_TIMESTAMP":utils.SixHoursAgo(),"COUNT":"3"}
    response,err := s.GetParamed("https://api.instagram.com/v1/users/"+id+"/media/recent/", params)
    if err != nil {
        panic(err)
    }

    //Decode request
    var posts Posts
    decoder := json.NewDecoder(response.Body)
    err = decoder.Decode(&posts)
    if err != nil {
        panic(err)
    }

    return posts
}

func GetPosts(s *session.Session, hashtag string) Posts{

    response,err := s.Get("https://api.instagram.com/v1/tags/" + hashtag +"/media/recent")
    if err != nil {
        panic(err)
    }

    //Decode request
    var posts Posts
    decoder := json.NewDecoder(response.Body)
    err = decoder.Decode(&posts)
    if err != nil {
        panic(err)
    }

    return posts
}

func GetUser(s *session.Session, id string) User{
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


func GetTag(s *session.Session, hashtag string) Tag{

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

func GetNext(s *session.Session, url string) Users{
    response,err := s.RawGet(url)
    if err != nil {
        panic(err)
    }

    //Decode request
    var bunch Users
    data, err := ioutil.ReadAll(response.Body)
    if err == nil && data != nil {
        err = json.Unmarshal(data, &bunch)
    }
    if err != nil {
        s.Log(string(data[:]))
        panic(err)
    }

    return bunch
}

func GetNextPost(s *session.Session, url string) Posts{
    response,err := s.RawGet(url)
    if err != nil {
        panic(err)
    }

    //Decode request
    var bunch Posts
    data, err := ioutil.ReadAll(response.Body)
    if err == nil && data != nil {
        err = json.Unmarshal(data, &bunch)
    }
    if err != nil {
        s.Log(string(data[:]))
        panic(err)
    }

    return bunch
}

func getPeople(s *session.Session, url string) (users Users){
    response,err := s.Get(url)
    if err != nil {
        panic(err)
    }

    data, err := ioutil.ReadAll(response.Body)
    if err == nil && data != nil {
        err = json.Unmarshal(data, &users)
    }
    if err != nil {
        s.Log(string(data[:]))
        panic(err)
    }

    return
}

func GetFollowing(s *session.Session) Users{
    return getPeople(s, "https://api.instagram.com/v1/users/" + s.GetId() +"/follows")
}

func GetFollowers(s *session.Session) Users{
    return getPeople(s, "https://api.instagram.com/v1/users/" + s.GetId() +"/followed-by")
}

func IsFollowing(s *session.Session, id string) bool {
    response ,err := s.Get("https://api.instagram.com/v1/users/"+id+"/relationship")
    if err != nil {
        panic(err)
    }

    var status Status
    decoder := json.NewDecoder(response.Body)
    err = decoder.Decode(&status)
    if err != nil {
        panic(err)
    }

    return status.Data.Outgoing_status == "follows"
}

// Post actions
func LikePosts(s *session.Session, id string) {
    v := url.Values{}

    response ,err := s.Post("https://api.instagram.com/v1/media/"+id+"/likes",v)
    if err != nil {
        panic(err)
    }
    s.Log(string(response.StatusCode))
}

func FollowUser(s *session.Session, id string){
    v := url.Values{}
    v.Set("action", "follow")

    response,err := s.Post("https://api.instagram.com/v1/users/"+ strings.Split(id,"_")[1] +"/relationship",v)
    if err != nil {
        panic(err)
    }
    s.Log(string(response.StatusCode))
}
