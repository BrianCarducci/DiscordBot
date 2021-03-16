#!/bin/bash

export BOT_DIR=/home/ubuntu/DiscordBot
export LOGS_DIR=$BOT_DIR/logs
export CGO_CFLAGS="-O2"

# chown/chgrp/chmod DiscordBot directory so ubuntu user can write logs
sudo chown -R ubuntu $BOT_DIR
sudo chgrp -R ubuntu $BOT_DIR
sudo chmod -R 744 $BOT_DIR

cd $BOT_DIR

/snap/bin/go build main.go >> $LOGS_DIR/after-install.log 2>&1
