#!/bin/bash

export BOT_DIR=/home/ubuntu/DiscordBot
export LOGS_DIR=$BOT_DIR/logs

export AWS_SDK_LOAD_NONDEFAULT_CONFIG="true"

cd $BOT_DIR

./main >> $LOGS_DIR/bot.log 2>&1 &
echo "PID of running DiscordBot: $!" >> $LOGS_DIR/bot.log 2>&1 &
