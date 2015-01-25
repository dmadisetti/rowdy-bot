package bot

type Settings struct {

    Errored bool
    Target float64
    Magic float64

    // App specifics
    Client_id string
    Client_secret string
    Callback string
    Hashtags []string

    // Account Specific
    Id string
    Access_token string
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

func (s *Settings) Valid() bool{
    return s.Id != "" && s.Client_id != "" && s.Client_secret != "" && s.Callback != ""
}