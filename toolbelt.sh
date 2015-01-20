#!/bin/sh

show_help(){
    echo "
    Why Hello there! You must be looking for help\n\
    \n\
    The Flags: \n\
    r - run \n\
    t - test \n\
    d - deploy \n\
    b - backup \n\
    i - init fom backup \n\
    s - setup\n\
    p - ci push
    c - clean
    \n\
    Chain em together as you see fit \n\
    "
}

APP_ID=rowdy-bot
current_datetime=$(date '+%Y%m%d_%H%M%S')
filename="backup_$current_datetime.txt"

setup(){
    export FILE=go_appengine_sdk_linux_amd64-1.9.17.zip
    curl -qO https://storage.googleapis.com/appengine-sdks/featured/$FILE
    unzip -q $FILE
}

run(){
    ./go_appengine/goapp serve;
}

try(){
    ./go_appengine/goapp build ./bot;
    ./go_appengine/goapp test ./tests;
}

deploy(){
    echo $PASSWORD | go_appengine/appcfg.py --email=$EMAIL --passin update ./
}

backup(){
    go_appengine/appcfg.py download_data --application=$APP_ID --url=http://$APP_ID.appspot.com/_ah/remote_api --filename=backups/$filename --email=$EMAIL;
}

init(){
    appcfg.py upload_data --application=$APP_ID --filename=backups/$filename --url=http://localhost:8080/_ah/remote_api --email=$EMAIL;
}

push(){
    try || exit 1;
    git branch | grep "\*\ [^(master)\]" || {
        deploy;
    }
}

clean(){
    rm -rf go_appengine*;
    rm bulkloader*;
}

while getopts "h?rtpsibcdx:" opt; do
    case "$opt" in
    h|\?)
        show_help
        ;;
    s)  setup
        ;;
    d)  deploy
        ;;
    b)  backup
        ;;
    i)  init
        ;;
    r)  run
        ;;
    t)  try
        ;;
    p)  push
        ;;
    c)  clean
        ;;
    esac
done
