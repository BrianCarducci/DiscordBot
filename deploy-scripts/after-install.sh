#!/bin/bash

export BOT_FOLDER=/home/ubuntu/DiscordBot
cd $BOT_FOLDER

go build main.go > after-install.log