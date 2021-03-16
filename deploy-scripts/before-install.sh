#!/bin/bash

export BOT_DIR=/home/ubuntu/DiscordBot
export LOGS_DIR=$BOT_DIR/logs

# Install gcc for gopus library
sudo apt-get -y install build-essential >> $LOGS_DIR/before-install.log 2>&1

# chown/chgrp/chmod DiscordBot directory so ubuntu user can write logs
sudo chown -R ubuntu $BOT_DIR
sudo chgrp -R ubuntu $BOT_DIR
sudo chmod -R 744 $BOT_DIR
