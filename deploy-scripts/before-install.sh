#!/bin/bash

export BOT_FOLDER=/home/ubuntu/DiscordBot

# rm directory if it exists
if [ -d $BOT_FOLDER ]
then
 rm -rf $BOT_FOLDER && echo "Removed pre-existing DiscordBot dir" > before-install.log || echo "Could not remove DiscordBot dir" > before-install.log
else
 echo "$BOT_FOLDER did not exist. Did nothing..." > before-install.log
fi

#mkdir -p $BOT-FOLDER