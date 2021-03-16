#!/bin/bash

export BOT_DIR=/home/ubuntu/DiscordBot
export LOGS_DIR=$BOT_DIR/logs

# Install necessary apt utilities
echo "Installing apt packages" >> $LOGS_DIR/before-install.log
{
cat <<-EOF
# Install gcc for gopus library
build-essential
EOF
} | grep -vE '^#' | xargs sudo apt-get -y install >> $LOGS_DIR/before-install.log 2>&1

# Install necessary snap utilities
printf "\nInstalling snap packages\n" >> $LOGS_DIR/before-install.log
{
cat <<-EOF
#--channel=latest/stable go
# Install ffmpeg to convert Polly audio sample rate to 48KHz
ffmpeg
EOF
} | grep -vE '^#' |
{
while IFS= read -r line; do
	echo "Attempting to install $line"
	sudo snap install "$line"
done >> $LOGS_DIR/before-install.log 2>&1
}
