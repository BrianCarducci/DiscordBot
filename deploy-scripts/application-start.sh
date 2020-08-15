#!/bin/bash

export BOT_DIR=/home/ubuntu/DiscordBot
export LOGS_DIR=$BOT_DIR/logs

cd $BOT_DIR

./main >> $LOGS_DIR/bot.log 2>&1 &