package bot

type Settings struct {

    Errored bool
    Target float64
    Magic float64

    // App specifics (Set these)
    Client_id string
    Client_secret string
    Callback string
    Hashtags []string

    // Account Specific
    Id string
    Access_token string

    // Bot Specific
    Hash string
    Production string
}

func NewSettings()*Settings{
    return &Settings{
        Errored : false,
        Target  : 1000,
        Magic   : 0.75,
        Client_id : "",
        Client_secret: "",
        Callback: "",
    }
}

func (s *Settings) Valid() bool{
    return s.Client_id != "" && s.Client_secret != "" && s.Callback != ""
}

func (s *Settings) Usable() bool{
    return s.Id != "" && s.Access_token != "" && len(s.Hashtags) > 0 
}