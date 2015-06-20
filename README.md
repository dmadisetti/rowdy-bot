Rowdy-Bot - Because our stuffed dog deserves to be famous
=========

[![Build Status](https://travis-ci.org/dmadisetti/rowdy-bot.png)](https://travis-ci.org/dmadisetti/rowdy-bot)

Put together over the weekend and improved off and on over a series of months. To read more about this project, [checkout this blog post](http://blog.postmodern.technology/machine-learning-instagram-bot)

**Note**:Our instagram client got banned. The Machine Learning part of this bot is a little bit to greedy with data. I will no longer be improving this repository, however, if you do in fact decide to fork `rowdy-bot`, you might want to turn down the number of requests- instead of pushing the bot to its limits.

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

    // Use so remote ML trainer can talk to production
    Hash string // Preshared key. Not super secure, but better than anyone messing with the service without minimal effort.
    Production string // Where do you want to share the data to?

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
Behind the Scenes
---

The bot uses standard [Round-Robin](https://en.wikipedia.org/wiki/Round-robin) for preset hashtags if Machine Learning is turned off.

If machine learning is turned on. The bot creates a model of (followers) vs (people the bot follows but don't follow back). Here's a graph of some of the data we came across:

![Data Points](https://raw.github.com/dmadisetti/rowdy-bot/master/FFP.png "Followers, Non-Followers and Posts")

Axis should be labeled Followers, Following, Posts (X,Y,Z) where Blue dots are the bot's followers and Red dots are the folk who don't follow the bot back.

Upon generating this model, we then use a permutted page rank for the hashtags the bot should follow.

The model generation and hashtag ranking was done on an old computer running Debian, just so we didn't rackup server fees (It's a heavy data load), and the processed data was pushed to production

---
Objectives:

- More Go - `Check`
- Become Instafamous - `Sorta?`
- Expand bot to hashtag crawl and use ML - `Of sorts`

Todo:

- Not get sued