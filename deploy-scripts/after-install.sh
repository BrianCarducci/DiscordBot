#!/bin/bash

export BOT_DIR=/home/ubuntu/DiscordBot
export LOGS_DIR=$BOT_DIR/logs
export CGO_CFLAGS="-O2"

cd $BOT_DIR

# Install gcc for gopus library
sudo apt-get -y install build-essential

/snap/bin/go build main.go >> $LOGS_DIR/after-install.log 2>&1
