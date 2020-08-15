#!/bin/bash

export BOT_DIR=/home/ubuntu/DiscordBot
export LOGS_DIR=$BOT_DIR/logs

cd $BOT_DIR

/snap/bin/go build main.go >> $LOGS_DIR/after-install.log 2>&1