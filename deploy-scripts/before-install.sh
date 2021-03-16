#!/bin/bash

export BOT_DIR=/home/ubuntu/DiscordBot
export LOGS_DIR=$BOT_DIR/logs

# Install necessary apt utilities
grep -vE '^#' $BOT_DIR/deploy-scripts/apt-utils.txt | xargs sudo apt-get -y install >> $LOGS_DIR/before-install.log 2>&1

# Install necessary snap utilities
while IFS= read -r line; do
	sudo snap install "$line" >> $LOGS_DIR/before-install.log 2>&1
done < <(grep -vE '^#' $BOT_DIR/deploy-scripts/snap-utils.txt)
