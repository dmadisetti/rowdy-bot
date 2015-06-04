package session

import "bot/utils"

type Machine struct {

    // Learning States
    Learning bool
    Learned bool

    // Next Url if job didn't finish
    Next string

    // Account Specific
    Step int

    // Calls (pretty much just cosemetic for current status)
    FollowingCalls int
    FollowerCalls int
    HashtagCalls int

    // Status
    Status int
    Processing bool

    // For hashtags
    BadSize int
    GoodSize int
    
    // Linear regression Params
    Lambda float64
    Alpha float64

    // Stored regression data
    Bias float64
    Xfollowers float64
    Xfollowing float64
    Xposts float64
}

func NewMachine() *Machine{
    return &Machine{
        Learning : false,
        Learned : false,
        Step : utils.APPRAISE,
        FollowingCalls : 0,
        FollowerCalls : 0,
        HashtagCalls : 0,
        Lambda : 10000,
        Alpha : 0.01,
        Status : 0,
        Bias : 0.0,
        Xfollowers : 0.0,
        Xfollowing : 0.0,
        Xposts : 0.0,
    }
}

func (m *Machine) SetLimits(following int, followers int){

    // Limits by calls
    m.FollowingCalls = int(following/utils.MAXREQUESTS) + 1
    m.FollowerCalls  = int(followers/utils.MAXREQUESTS) + 1
    m.HashtagCalls   = m.FollowerCalls + m.FollowerCalls

}

func (m *Machine) GetLimit() int {
    switch m.Step {
    case utils.APPRAISE:
        return m.FollowerCalls
    case utils.SCORN:
        return m.FollowingCalls
    default:
        return m.HashtagCalls
    }
}
