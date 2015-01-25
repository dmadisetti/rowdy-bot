package bot

import (
    "net/http"
    "appengine"
    "appengine/urlfetch"
    "appengine/datastore"
    "net/url"
    "encoding/json"
    "bytes"
)

type Session struct {
    settings *Settings
    context appengine.Context
    client *http.Client
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
        settings: NewSettings(),
    }
    return s
}

// Talk with the data layer
func (session *Session) Load() bool{
    err := datastore.Get(session.context,datastore.NewKey(session.context,"Settings","",1, nil),session.settings)
    return err != nil || !session.Valid()
}

func (session *Session) Save(){
    datastore.Put(session.context,datastore.NewKey(session.context,"Settings","",1, nil),session.settings)
}

func (session *Session) Valid() bool{
    return session.settings.Valid()
}

// HTTP functions
func (session *Session) Get(uri string) (*http.Response, error){
    request,err := http.NewRequest("GET", uri +"?client_id=" + session.settings.Client_id, nil)
    if err != nil {
        panic(err)
    }
    session.Sign(*request)
    return session.client.Do(request)
}

func (session *Session) Post(uri string, v url.Values) (*http.Response, error){
    session.Authenticate(v)
    request,err := http.NewRequest("POST", uri, bytes.NewBufferString(v.Encode()))
    if err != nil {
        panic(err)
    }
    session.Sign(*request)
    return session.client.Do(request)
}

// Might be better breaking into actions
func (session *Session) SetAuth(code string){

    v := url.Values{}
    v.Set("client_id",session.settings.Client_id)
    v.Add("client_secret",session.settings.Client_secret)
    v.Add("grant_type","authorization_code")
    v.Add("redirect_uri",session.settings.Callback)
    v.Add("code",code)

    request,err := http.NewRequest("POST", "https://api.instagram.com/oauth/access_token", bytes.NewBufferString(v.Encode()))
    if err != nil {
        panic(err)
    }

    session.Sign(*request)
    response,err := session.client.Do(request)

    //Decode request
    var auth Auth
    decoder := json.NewDecoder(response.Body)
    err = decoder.Decode(&auth)
    if err != nil {
        panic(err)
    }

    session.context.Infof("Interval: %v",request)

    session.settings.Access_token = auth.Access_token
    session.settings.Id = auth.User.Id
    session.Save()
}


// Hashtags!
func (session *Session) HasHashtags() bool{
    return len(session.settings.Hashtags) > 0
}

func (session *Session) SetHashtags(tags []string){
    session.settings.Hashtags = tags
    session.Save()
}

// Getters
func (s *Session) GetHashtag(intervals int) (hashtag string){
    hashtag = s.settings.Hashtags[intervals % len(s.settings.Hashtags)]
    // Some logging
    s.context.Infof("Hashtag: %v",hashtag)
    s.context.Infof("Interval: %v",intervals)
    return
}
func (s *Session) GetId() string{
    return s.settings.Id
}
func (s *Session) GetMagic() float64 {
    return s.settings.Magic
}
func (s *Session) GetTarget() float64 {
    return s.settings.Target
}

// For rendering
func (s *Session) GetAuthLink() string{
    return "https://instagram.com/oauth/authorize/?client_id="+ s.settings.Client_id + "&response_type=code&scope=likes+comments+relationships&redirect_uri=" + s.settings.Callback
}

func (s *Session) GetHashtags() []string{
    return s.settings.Hashtags
}

// Http helpers
func (s *Session) Authenticate(v url.Values){
    v.Add("access_token", s.settings.Access_token)
}

func (s *Session) Sign(request http.Request){
    ip := "127.0.0.1"
    request.Header.Set("X-Insta-Forwarded-For", ip + "|" + ComputeHmac256(ip, s.settings.Client_secret))
}
