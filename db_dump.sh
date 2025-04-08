#!/usr/bin/bash

cd /home/vedomir/go/src/remindista
git pull
/home/vedomir/bin/task postgres:dump
git add ./db/create_tables.sql && git commit -m "postgres dump $(date)" && git push
