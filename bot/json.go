package bot

type Counts struct {
    Follows int64
    Followed_by int64
}

type Data struct {
    Id string
    Counts Counts
    Media_count int64
    Name string
}

type Status struct {
    Data Data
}

type User struct {
    Id string
    Data Data
}

type Auth struct {
    Access_token string // Can parse apart for user
    User User
}

type Post struct {
    Id string // Can parse apart for user
}

type Posts struct {
    Data []Post
}

type Tag struct {
    Data Data
}
