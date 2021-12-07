#!/bin/sh

find /root/views/js -name 'main.js' | xargs sed -i "s http://localhost:8080 $API_HOST g"
