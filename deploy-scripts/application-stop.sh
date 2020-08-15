#!/bin/bash

set -e

export BOT_DIR=/home/ubuntu/DiscordBot
export LOGS_DIR=$BOT_DIR/logs

{
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
  PROCS=$(ps -ax | grep './main' | grep -v 'grep' | awk '{ print $1 }')
  if [ -z "$PROCS" ]
  then
    echo "No running DiscordBots."
  else
    echo "Killing the following processes:"
    ps -p $PROCS
    kill -9 $PROCS
  fi
} >> $LOGS_DIR/application-stop.log 2>&1 