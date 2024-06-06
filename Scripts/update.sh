#!/bin/bash
cd /home/sammy/server
./updater > /home/sammy/Updater.log
./autoFilterAll > /home/sammy/Filter.log &
cd /home/sammy/Test/www
./makePromXlsx.php > /home/sammy/Prom.log &
