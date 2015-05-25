package bot

import(
    "log"
    "strings"
    "io/ioutil"
    "net/url"
    "encoding/json"
)

// Could dynamically build these, but naw

// Get actions
func GetStatus(s *Session) (count Counts){

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

func GetMedia(s *Session, id string) Posts{
    params := map[string]string{"MIN_TIMESTAMP":SixHoursAgo(),"COUNT":"3"}
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

func GetPosts(s *Session, hashtag string) Posts{

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

    return posts
}

func GetUser(s *Session, id string) User{
    log.Println(id)
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

func GetNext(s *Session, url string) Users{
    log.Println(url)
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
        log.Println(string(data[:]))
        panic(err)
    }

    return bunch
}

func getPeople(s *Session, url string) (users Users){
    response,err := s.Get(url)
    if err != nil {
        panic(err)
    }

    data, err := ioutil.ReadAll(response.Body)
    if err == nil && data != nil {
        err = json.Unmarshal(data, &users)
    }
    if err != nil {
        log.Println(string(data[:]))
        panic(err)
    }

    return
}

func GetFollowing(s *Session) Users{
    return getPeople(s, "https://api.instagram.com/v1/users/" + s.GetId() +"/follows")
}

func GetFollowers(s *Session) Users{
    return getPeople(s, "https://api.instagram.com/v1/users/" + s.GetId() +"/followed-by")
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
