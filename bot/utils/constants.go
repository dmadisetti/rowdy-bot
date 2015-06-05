package utils

// Time constants
const DAY      int64   = 60 * 60 * 24
const INTERVAL float64 = 60 * 5 // In seconds. Should match Cron Job
const HOUR     float64 = 60 * 60
const SIXHOURS float64 = HOUR * 6

// Process Limits
const FOLLOWS       int = 60
const LIKES         int = 100
const MAX           int = 5000
const GRABCOUNT     int = 50
const MAXPOSTGRAB   int = 4
const CALLS         int = int(HOUR/INTERVAL)
const MAXREQUESTS   int = int(MAX/CALLS)
const MAXPEOPLEGRAB int = int(MAXREQUESTS/GRABCOUNT)

// Process Steps
const APPRAISE    int = 0
const SCORN       int = 1
const BUILD       int = 2
const GOODTAGS    int = 3
const BADTAGS     int = 4
const COMPUTETAGS int = 5
const SHARE       int = 6
