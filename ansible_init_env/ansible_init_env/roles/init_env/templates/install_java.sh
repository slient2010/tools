#!/bin/bash

[ ! -d "{{ app_dir }}/{{ java_version }}" ] && exit 
ln -s {{ app_dir }}/{{ java_version }}/bin/java /usr/local/bin
echo "
export JAVA_HOME={{ app_dir }}/{{ java_version }}
export PATH={{ app_dir }}/{{ java_version }}/bin:$PATH 
export CLASSPATH=.:{{ app_dir }}/{{ java_version }}/lib/dt.jar:{{ app_dir }}/{{ java_version }}/lib/tools.jar 
" >> /etc/profile
rm -rf /tmp/$(basename $0)