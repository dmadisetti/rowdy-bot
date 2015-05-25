package bot

// Process Steps
const APPRAISE int = 0
const SCORN int = 1
const BUILD int = 2
const GOODTAGS int = 3
const BADTAGS int = 4
const COMPUTETAGS int = 5
const SHARE int = 6

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
        Step : APPRAISE,
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
    m.FollowingCalls = int(following/MAXREQUESTS) + 1
    m.FollowerCalls  = int(followers/MAXREQUESTS) + 1
    m.HashtagCalls   = m.FollowerCalls + m.FollowerCalls

}

func (m *Machine) GetLimit() int {
    switch m.Step {
    case APPRAISE:
        return m.FollowerCalls
    case SCORN:
        return m.FollowingCalls
    default:
        return m.HashtagCalls
    }
}
