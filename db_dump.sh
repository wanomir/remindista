#!/usr/bin/bash

cd /home/vedomir/go/src/remindista
git pull
/snap/bin/task postgres:dump
git add ./db/create_tables.sql && git commit -m "postgres dump $(date)" && git push
