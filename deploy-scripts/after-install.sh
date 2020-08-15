#!/bin/bash

export BOT_DIR=/home/ubuntu/DiscordBot
cd $BOT_DIR

go build main.go > after-install.log