#!/bin/bash

cd /tmp/{{ ngx_openresty_version }}
./configure --prefix={{ app_dir }}/openresty
make > /dev/null 2>&1
make install > /dev/null 2>&1
 
echo "
export PATH=$PATH:{{ app_dir }}/openresty/nginx/sbin
" >> /etc/profile

touch /tmp/asdfasdfasfdad
rm -rf /tmp/$(basename $0)
