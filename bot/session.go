package bot

import (
    "net/http"
    "appengine"
    "appengine/urlfetch"
    "appengine/datastore"
    "net/url"
    "bytes"
    "log"
    "encoding/json"
    "crypto/hmac"
    "crypto/sha256"
    "encoding/hex"
)

type Session struct {
    Settings *Settings
    context appengine.Context
    client *http.Client
}

type Settings struct {

    Errored bool
    Target float64
    Magic float64

    Id string
    Client_id string
    Client_secret string
    Callback string
    Hashtags []string

    Access_token string
}


func NewSession(r *http.Request) *Session {

    c := appengine.NewContext(r)
    cl := urlfetch.Client(c)

    // Set transport to allow for https
    cl.Transport = &urlfetch.Transport{
        Context:                       c,
        Deadline:                      0,
        AllowInvalidServerCertificate: false,
    }

    s := &Session{
        context : c,
        client  : cl,
        Settings: NewSettings(),
    }
    return s
}

func (session *Session) Load() bool{
    err := datastore.Get(session.context,datastore.NewKey(session.context,"Settings","",1, nil),session.Settings)
    return err != nil || !session.Valid()
}

func (session *Session) Save(){
    datastore.Put(session.context,datastore.NewKey(session.context,"Settings","",1, nil),session.Settings)
}

func (session *Session) Valid() bool{
    return session.Settings.Valid()
}

func (session *Session) Get(uri string) (*http.Response, error){
    log.Println(uri)
    request,err := http.NewRequest("GET", uri +"?client_id=" + session.Settings.Client_id, nil)
    if err != nil {
        panic(err)
    }
    session.Settings.Sign(*request)
    return session.client.Do(request)
}

func (session *Session) Post(uri string, v url.Values) (*http.Response, error){
    session.Settings.Authenticate(v)
    request,err := http.NewRequest("POST", uri, bytes.NewBufferString(v.Encode()))
    if err != nil {
        panic(err)
    }
    session.Settings.Sign(*request)
    return session.client.Do(request)
}

// Might be better under actions
func (session *Session) SetAuth(code string){

    v := url.Values{}
    v.Set("client_id",session.Settings.Client_id)
    v.Add("client_secret",session.Settings.Client_secret)
    v.Add("grant_type","authorization_code")
    v.Add("redirect_uri",session.Settings.Callback)
    v.Add("code",code)

    request,err := http.NewRequest("POST", "https://api.instagram.com/oauth/access_token", bytes.NewBufferString(v.Encode()))
    if err != nil {
        panic(err)
    }

    session.Settings.Sign(*request)
    response,err := session.client.Do(request)
    log.Println(response)

    //Decode request
    var auth Auth
    decoder := json.NewDecoder(response.Body)
    err = decoder.Decode(&auth)
    if err != nil {
        panic(err)
    }

    session.Settings.Access_token = auth.Access_token
    session.Save()
}

func (session *Session) SetHashtags(tags []string){
    session.Settings.Hashtags = tags
    session.Save()
}

func (s *Settings) Valid() bool{
    return s.Id != "" && s.Client_id != "" && s.Client_secret != "" && s.Callback != ""
}

func (s *Settings) Authenticate(v url.Values){
    v.Add("access_token", s.Access_token)
}

func (s *Settings) Sign(request http.Request){
    ip := "127.0.0.1"
    request.Header.Set("X-Insta-Forwarded-For", ip + "|" + ComputeHmac256(ip, s.Client_secret))
}

func ComputeHmac256(message string, secret string) string {
    key := []byte(secret)
    h := hmac.New(sha256.New, key)
    h.Write([]byte(message))
    return hex.EncodeToString(h.Sum(nil))
}

func (s *Settings) GetHashtag(intervals float64) string{
    return s.Hashtags[int(intervals) % len(s.Hashtags)]
}

func (s *Settings) GetId() string{
    return s.Id
}

func NewSettings()*Settings{   
    return &Settings{
        Errored : false,
        Target  : 1000,
        Magic   : 0.75,
        Id      : "",
        Client_id : "",
        Client_secret: "",
        Callback: "",
    }
}