#!/bin/bash

# kill running bot
kill -9 $(ps -ax | grep "./main" | head -n 1 | awk '{ print $1 }') && echo "Kill succeeded!" > application-stop.log || echo "Kill FAILED!" > application-stop.log 