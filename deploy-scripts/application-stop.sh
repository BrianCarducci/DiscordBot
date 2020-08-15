#!/bin/bash

export BOT_DIR=/home/ubuntu/DiscordBot
export LOGS_DIR=$BOT_DIR/logs

# rm DiscordBot directory contents if they exist
if [ -d $BOT_DIR ]
then
 rm -rf $BOT_DIR/*
else
 mkdir -p $BOT_DIR
fi

# create logs directory
mkdir -p $BOT_DIR

# kill running bot
kill -9 $(ps -ax | grep "./main" | head -n 1 | awk '{ print $1 }') && echo "Kill succeeded!" >> $LOGS_DIR/application-stop.log 2>&1 || echo "Kill FAILED!" >> $LOGS_DIR/application-stop.log 2>&1 