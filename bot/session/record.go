package session

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
        s.context.Infof("Records: %v",Saved)
        s.SaveRecords(records)
        return
    }
    s.context.Infof("Error: %v", err)
    return
}

func (s *Session) SetRecords(follows, followed_by int64) {

    now := time.Now()

    t := strconv.FormatInt(now.Unix(),10)
    x := strconv.FormatInt(follows,10)
    y := strconv.FormatInt(followed_by,10)

    records := s.GetRecords()
    records.String += ",[" + t + ","+ x +","+ y +"]"

    s.SaveRecords(records)
}

func (s *Session) SaveRecords(records *Record){
    datastore.Put(s.context,datastore.NewKey(s.context,"Records","",1, nil),records)
}
