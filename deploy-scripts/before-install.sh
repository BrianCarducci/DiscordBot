#!/bin/bash

export BOT_DIR=/home/ubuntu/DiscordBot

# rm directory if it exists
if [ -d $BOT_DIR ]
then
 rm -rf $BOT_DIR && echo "Removed pre-existing DiscordBot dir" > before-install.log || echo "Could not remove DiscordBot dir" > before-install.log
else
 echo "$BOT_DIR did not exist. Did nothing..." > before-install.log
fi

mkdir -p $BOT_DIR