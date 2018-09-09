#!/bin/bash

ls /usr/local/bin/ | grep "redis-" > /dev/null 2>&1
[ $? -eq 0 ] && echo "redis already installed." && exit
cd {{ app_dir }}/{{ redis_version }}/
make 
make install
cd src
cp -rp `ls |grep redis- |grep -v "\."` /usr/local/bin/
# updatedb
rm -rf /tmp/$(basename $0)