package http

type Counts struct {
    Follows int64
    Followed_by int64
    Media int64
}

type Pagination struct{
    Next_url string // https:\/\/api.instagram.com\/v1\/users\/3\/media\/recent?access_token=184046392.f59def8.c5726b469ad2462f85c7cea5f72083c0&max_id=205140190233104928_3
}

type Data struct {
    Id string
    Counts Counts
    Media_count int64
    Username string
    Outgoing_status string
}

type Status struct {
    Data Data
}

type User struct {
    Id string
    Username string
    Data Data
}

type Users struct {
    Data []User
    Pagination Pagination
}

type Auth struct {
    Access_token string // Can parse apart for user
    User User
}

type Post struct {
    Id string // Can parse apart for user
    User User
    Tags []string
}

type Posts struct {
    Data []Post
    Pagination Pagination
}

type Tag struct {
    Data Data
}
