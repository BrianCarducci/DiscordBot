#!/bin/bash

export BOT_DIR=/home/ubuntu/DiscordBot
cd $BOT_DIR

/snap/bin/go build main.go > after-install.log