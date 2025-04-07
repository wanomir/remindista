#!/usr/bin/bash

cd /home/wanomir/go/src/remindista
git pull
/home/wanomir/bin/task postgres:dump
git add ./db/create_tables.sql && git commit -m "postgres dump $(date)" && git push
