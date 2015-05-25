Rowdy-Bot - Because our stuffed dog deserves to be famous
=========

[![Build Status](https://travis-ci.org/dmadisetti/rowdy-bot.png)](https://travis-ci.org/dmadisetti/rowdy-bot)

Put together over the weekend.

Our stuffed dog is awesome and deserves to become instagram famous. So I wrote him a bot. We're collecting follower data, so we'll see how it goes.

---
Manage upon hitting the site or from the GAE datastore viewer

```
	// Variables for exponential decay of follower ratio
    Target float64 
    Magic float64

    // App specifics (Set these upon visiting the page)
    Client_id string // Client in instagram app settings
    Client_secret string // Secret in instagram app settings
    Callback string // Callback configured in instagram app settings

    // Account Specific (Don't set these)
    Id string // Who ares you?
    Access_token string // The secret password

```

To run the application checkout the toolbelt
-----
From `./toolbelt -h`:
```
    Why Hello there! You must be looking for help
    
    The Flags: 
    r - run 
    t - test 
    d - deploy 
    b - backup 
    i - init fom backup 
    s - setup
    l - train model
    p - ci push
    c - clean
    
    Chain em together as you see fit 
```

To get started run `./toolbelt -sr`


To set hashtags use the web interface
-----
![Our simple web interface](https://raw.github.com/dmadisetti/rowdy-bot/master/example.png "Screenshot")

---
This is our dog
-----
![Glorious pictures of Rowdy](https://raw.github.com/dmadisetti/rowdy-bot/master/rowdy.png "Screenshot of IG")

---
Objectives:

- More Go - `Check`
- Become Instafamous - `Sorta?`
- Expand bot to hashtag crawl and use ML - `In progess`

Todo:

- Follow back (Maybe)