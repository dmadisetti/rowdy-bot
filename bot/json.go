package bot

type Counts struct {
    Follows int64
    Followed_by int64
}

type Data struct {
    Counts Counts
}

type Status struct {
    Data Data
}

type Auth struct {
    Access_token string // Can parse apart for user
}

type Post struct {
    Id string // Can parse apart for user
}

type Posts struct {
    Data []Post
}