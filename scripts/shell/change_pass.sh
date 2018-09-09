#!/bin/bash

ip=`ifconfig eth0 |grep "inet addr"| cut -f 2 -d ":"|cut -f 1 -d " "`
vmmac=$(echo $ip|awk -F'.' '{printf("%02X%02X",$3,$4)}')
passwd="Mkk2015@${vmmac}321"
echo ${passwd} | passwd --stdin root
