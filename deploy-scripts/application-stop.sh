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
mkdir -p $LOGS_DIR

# kill running bot
PROCS="$(ps -ax | grep './main' | grep -v 'grep' | awk '{ print $1 }')"
if [ "$PROCS" | wc -l) -eq 0 ]
then
  echo "No running DiscordBots." >> $LOGS_DIR/application-stop.log 2>&1
else
  echo "Killing the following processes:" >> $LOGS_DIR/application-stop.log
  ps -p $PROCS >> $LOGS_DIR/application-stop.log 2>&1 
  kill -9 $PROCS && echo "Kill succeeded!" >> $LOGS_DIR/application-stop.log 2>&1 || echo "Kill FAILED!" >> $LOGS_DIR/application-stop.log 2>&1 
fi