#!/bin/bash

export BOT_DIR=/home/ubuntu/DiscordBot
export LOGS_DIR=$BOT_DIR/logs

# Install necessary apt utilities
echo "Installing apt packages" >> $LOGS_DIR/before-install.log
grep -vE '^#' $BOT_DIR/deploy-scripts/apt-utils.txt | xargs sudo apt-get -y install >> $LOGS_DIR/before-install.log 2>&1

# Install necessary snap utilities
printf "\nInstalling snap packages\n" >> $LOGS_DIR/before-install.log
while IFS= read -r line; do
	echo "Attempting to install $line"
	sudo snap install "$line"
done < <(grep -vE '^#' $BOT_DIR/deploy-scripts/snap-utils.txt) >> $LOGS_DIR/before-install.log 2>&1
