#!/bin/bash
if pgrep GetBets.php > /dev/null
then
  echo "GetBets.php is running"
else
  cd /home/sammy/SportPhp
  ./GetBets.php
fi
