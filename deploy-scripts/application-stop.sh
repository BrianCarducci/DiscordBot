#!/bin/bash

# kill running bot
kill -9 $(ps -ax | grep "./main" | head -n 1 | awk '{ print $1 }') && echo "Kill succeeded!" > application-stop.log || echo "Kill FAILED!" > application-stop.log 

export BOT_DIR=/home/ubuntu/DiscordBot

# rm directory if it exists
if [ -d $BOT_DIR ]
then
 rm -rf $BOT_DIR/* && echo "Removed everything in pre-existing DiscordBot dir ($BOT_DIR)" >> application-stop.log || echo "Could not remove files in DiscordBot dir" > before-install.log
else
 echo "$BOT_DIR did not exist. Creating it..." >> application-stop.log
 mkdir -p $BOT_DIR
fi
