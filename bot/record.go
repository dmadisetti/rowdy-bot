package bot

import (
    "appengine/datastore"
    "time"
    "strconv"
)

import(
)

// Record to save Followers/Followings on a daily basis
// Concat in String to save read write limit
type Record struct {
    String string
}

func (s *Session) GetRecords() (records *Record){
    records = &Record{}
    err := datastore.Get(s.context,datastore.NewKey(s.context,"Records","",1, nil),records)
    if err !=nil{
        s.SaveRecords(records)
    }
    return
}

func (s *Session) SetRecords(count Counts) {

    now := time.Now()

    t := strconv.FormatInt(now.Unix(),10)
    x := strconv.FormatInt(count.Follows,10)
    y := strconv.FormatInt(count.Followed_by,10)

    records := s.GetRecords()
    records.String += ",[" + t + ","+ x +","+ y +"]"

    s.SaveRecords(records)
}

func (s *Session) SaveRecords(records *Record){
    datastore.Put(s.context,datastore.NewKey(s.context,"Records","",1, nil),records)
}
