package session

// Session getting big. 
// Might be a good idea
// To break out ML
// - A haiku

import (
    "appengine"
    "appengine/datastore"
    "appengine/urlfetch"
    "appengine/memcache"
    "bytes"
    "encoding/json"
    "net/http"
    "net/url"
    "strings"
    "bot/utils"
)

type Session struct {
    settings *Settings
    machine *Machine
    context appengine.Context
    client *http.Client

    // Particular to machine learning sessions
    keys []*datastore.Key
    people []Person
    count int
    processed int
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
        machine: NewMachine(),
    }
    return s
}

// Talk with the data layer
func (session *Session) LoadSettings() bool{
    err := datastore.Get(session.context,datastore.NewKey(session.context,"Settings","",1, nil),session.settings)
    return !(err != nil || !session.Valid())
}
func (session *Session) LoadMachine() bool{
    err := datastore.Get(session.context,datastore.NewKey(session.context,"Machine","",1, nil),session.machine)
    return !(err != nil || !session.Valid())
}

func (session *Session) GetPeople() (people []Person){
    datastore.NewQuery("Person").GetAll(session.context, &people)
    return
}

func (session *Session) GetPeopleCursor(positive bool, offset int) *datastore.Iterator{
    return datastore.NewQuery("Person").Filter("Follows = ", positive).Limit(400).Offset(offset).KeysOnly().Run(session.context)
}

func (session *Session) GetHashtagCursor() *datastore.Iterator{
    return datastore.NewQuery("Hashtag").KeysOnly().Run(session.context)
}

func (session *Session) SetTopTags(){
    var hashtags []Hashtag
    datastore.NewQuery("Hashtag").Limit(utils.CALLS).Order("-Value").GetAll(session.context, &hashtags)
    var tags []string
    for _, hashtag := range hashtags {
        tags = append(tags,hashtag.Name)
    }
    session.SetHashtags(tags)
}

func (session *Session) Hashtag(tag string) (hashtag Hashtag){
    datastore.Get(session.context, session.key("Hashtag",tag), &hashtag)
    hashtag.Name = tag
    return
}

// Save
func (session *Session) Save(){
    session.SaveSettings()
    session.SaveMachine()
}
func (session *Session) SaveSettings(){
    datastore.Put(session.context,datastore.NewKey(session.context,"Settings","",1, nil) ,session.settings)
}
func (session *Session) SaveMachine(){
    datastore.Put(session.context,datastore.NewKey(session.context,"Machine" ,"",1, nil) ,session.machine)
}
func (session *Session) SavePeople(){
    datastore.PutMulti(session.context, session.keys, session.people)
}
func (session *Session) SaveHashtag(hashtag Hashtag){
    datastore.Put(session.context, session.key("Hashtag",hashtag.Name), &hashtag)
}
func (session *Session) key(entity, id string) (*datastore.Key){
    return datastore.NewKey(session.context, entity, id, 0, nil)
}
func (session *Session) PutPerson(person Person, id string){
    session.keys   = append(session.keys,session.key("Person", id))
    session.people = append(session.people,person)
}

// Flush
func (session *Session) Flush(){
    session.FlushEntity("Person")
    session.FlushEntity("Hashtag")

    memcache.Flush(session.context)

    session.machine = NewMachine()
    session.SaveMachine()
}
func (session *Session) FlushEntity(entity string){
    keys, _ := datastore.NewQuery(entity).KeysOnly().GetAll(session.context,nil)
    datastore.DeleteMulti(session.context, keys)
}

// Cache
func (session *Session) CheckCache(id string) bool{
    item := &memcache.Item{
        Key:   id,
        Value: []byte(""),
    }
    // Add the item to the memcache, if the key does not already exist
    if err := memcache.Add(session.context, item); err == memcache.ErrNotStored {
        return true
    }
    return false
}

// Session helpers
func (session *Session) IncrementCount(){
    session.count += 1
    return
}

func (session *Session) FinishedCount() bool{
    session.processed += 1
    return session.processed >= session.count
}

func (session *Session) CheckCount() bool{
    return session.processed >= session.count
}

func (session *Session) Share(){
    url := session.settings.Production 
    url += "/update"
    url += "?hash=" + session.settings.Hash
    url += "&hashtags="
    for _,tag := range session.settings.Hashtags {
        url += tag + ","
    }
    url += utils.FloatToString(session.machine.Bias) + " "
    url += utils.FloatToString(session.machine.Xfollowing) + " "
    url += utils.FloatToString(session.machine.Xfollowers) + " "
    url += utils.FloatToString(session.machine.Xposts) + " "
    session.RawGet(url)
}

func (session *Session) VerifiedUpdate(secret string) bool{
    return secret == session.settings.Hash
}

// Settings Helpers
func (session *Session) Valid() bool{
    return session.settings.Valid()
}

func (session *Session) InitAuth(client_id, client_secret, callback, hash string){
    if !session.Valid(){
        session.settings.Client_id = client_id
        session.settings.Client_secret = client_secret
        session.settings.Callback = callback
        session.settings.Hash = hash
        session.SaveSettings()
    }
}

// Machine helpers
func (s *Session) Warn(code int){
    s.context.Warningf("Bad Code: %v",code)
}
func (session *Session) IncrementStep(){
    session.machine.Step += 1
    session.machine.Status = 0
    session.SaveMachine()
    return
}
func (session *Session) IncrementState() float64{
    session.machine.Status += 1
    session.SaveMachine()
    return float64(session.machine.Step)
}
func (session *Session) IncrementSize(size int, positive bool) {
    if positive {
        session.machine.GoodSize += size
    }else{
        session.machine.BadSize += size
    }
    session.SaveMachine()    
}

// HTTP functions
func (session *Session) Get(uri string) (*http.Response, error){
    request,err := http.NewRequest("GET", uri +"?access_token="+ session.settings.Access_token +"&client_id=" + session.settings.Client_id, nil)
    if err != nil {
        panic(err)
    }
    session.Sign(*request)
    return session.client.Do(request)
}

func (session *Session) GetParamed(uri string, params map[string]string) (*http.Response, error){
    uri += "?client_id=" + session.settings.Client_id
    for key, value := range params {
        uri += "&" + key + "=" + value
    }
    request,err := http.NewRequest("GET", uri, nil)
    if err != nil {
        panic(err)
    }
    session.Sign(*request)
    return session.client.Do(request)
}


func (session *Session) RawGet(uri string) (*http.Response, error){
    request,err := http.NewRequest("GET", uri, nil)
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
func (session *Session) Auth(code string) *json.Decoder{

    v := url.Values{}
    v.Set("client_id",session.settings.Client_id)
    v.Add("client_secret",session.settings.Client_secret)
    v.Add("grant_type","authorization_code")
    v.Add("code",code)
    v.Add("redirect_uri",session.settings.Callback)

    // Hack to prevent ?+ getting encoded
    data := v.Encode() + "?hashtags=" + strings.Join(session.settings.Hashtags,"+")

    request,err := http.NewRequest("POST", "https://api.instagram.com/oauth/access_token", bytes.NewBufferString(data))
    if err != nil {
        panic(err)
    }

    session.Sign(*request)
    response,err := session.client.Do(request)
    if err != nil {
        panic(err)
    }

    decoder := json.NewDecoder(response.Body)
    return decoder
}

// Might be better breaking into actions
func (session *Session) SetAuth(token, id string){
    session.settings.Access_token = token
    session.settings.Id = id
    session.SaveSettings()
}


// Basic decision Hashtags!
func (session *Session) Usable() bool{
    return session.settings.Usable()
}
func (session *Session) SetHashtags(tags []string){
    session.settings.Hashtags = tags
    session.SaveSettings()
}

// Setters
// Machine
func (session *Session) SetLearning() bool{
    set := session.machine.Learning
    session.machine.Learning = true
    session.SaveMachine()
    return !set
}
func (session *Session) SetLearnt(){
    session.machine.Learned = true
    session.SaveMachine()
}
func (session *Session) SetLimits(followers, following int){
    session.machine.SetLimits(followers, following)
    session.SaveMachine()
}
func (session *Session) SetNext(next string){
    session.machine.Next = next;
    session.SaveMachine()
}
func (session *Session) ParseTheta(theta []string){
    session.SetTheta([]float64{
        utils.StringToFloat(theta[0]),
        utils.StringToFloat(theta[1]),
        utils.StringToFloat(theta[2]),
        utils.StringToFloat(theta[3]),
    })
}
func (session *Session) SetTheta(theta []float64){
    session.machine.Bias       = theta[0]
    session.machine.Xfollowers = theta[1]
    session.machine.Xfollowing = theta[2]
    session.machine.Xposts     = theta[3]
    session.SaveMachine()
}

// Getters
// Settings
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

// Getters
// Machine
func (session *Session) GetLimit() int{
    return session.machine.GetLimit()
}
func (session *Session) GetLearnt() bool{
    return session.machine.Learned
}
func (session *Session) GetLearningStep() int{
    return session.machine.Step
}
func (session *Session) GetNext() string{
    return session.machine.Next
}
func (session *Session) GetState() int{
    return session.machine.Status
}
func (session *Session) GetLambda() float64{
    return session.machine.Lambda
}
func (session *Session) GetAlpha() float64{
    return session.machine.Alpha
}
func (session *Session) GetHashtagSize(positive bool) float64{
    if positive {
        return float64(session.machine.GoodSize)
    }
    return float64(session.machine.BadSize)
}
func (session *Session) GetTheta() []float64{
    return []float64{
        -session.machine.Bias,
        -session.machine.Xfollowers,
        -session.machine.Xfollowing,
        -session.machine.Xposts,
    }
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
    request.Header.Set("X-Insta-Forwarded-For", ip + "|" + utils.ComputeHmac256(ip, s.settings.Client_secret))
}
