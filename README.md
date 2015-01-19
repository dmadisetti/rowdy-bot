Rowdy-Bot - Because our stuffed dog deserves to be famous
=========

[![Build Status](https://travis-ci.org/dmadisetti/rowdy-bot.png)](https://travis-ci.org/dmadisetti/rowdy-bot)

Put together over the weekend.

Our stuffed dog is awesome and deserves to become instagram famous. So I wrote him a bot. We're collecting follower data, so we'll see how it goes.

---
Manage settings from the GAE datastore viewer

```
	// Variables for exponential decay of follower ratio
    Target float64 
    Magic float64

    // Instgram Stuff
    Id string // Instagram ID of the user. (Can find this using source or reading through api docs)
    Client_id string // Client in instagram app settings
    Client_secret string // Secret in instagram app settings
    Callback string // Callback configured in instagram app settings

    Access_token string // Don't set this manually. Set when authorized

```

To set hash tags `curl` or just hit in your browser `domain.tld/hashtag?hashtags=comma,seperated,hashtags`

---
This is our dog
![alt tag](https://raw.github.com/dmadisetti/rowdy-bot/master/rowdy.png "Screenshot")

---
Objectives:

- More Go - `Check`
- Become Instafamous - `Pending`

Todo:

- Custom settings interface
- Easier hashtag settings
