#!/bin/bash
#
# This shell is used for manage mass machines
# 

current=`date +%Y_%m_%d`
[ ! -n "$1" ] && echo "Usage: $0 Command" && exit -1

#得到实验室物理机机器IP数组
source ip.sh
for id in ${!ip_list[@]}
do
    ip[$id]=${ip_list[$id]}
done

#echo ${!ip[@]}

#多线程
thread=50
tmp_fifofile="/tmp/$$.fifo"
mkfifo $tmp_fifofile
exec 6<>$tmp_fifofile
rm $tmp_fifofile

for ((i=0; i<$thread; i++))
do
    echo
done >&6


count=0

#获取对应服务器IP,执行操作

for j in ${!ip[@]};
do
    mip=${ip[$j]}
#    echo $mip
    ping -c 2 ${mip}>/dev/null 2>&1
    if [ $? -ne 0 ];
    then
        continue;
    fi
    read -u6
    {
        ./auto_action.sh ${mip} "$1" >./log/${mip}_${current}.log && {
           echo "======================${mip} task already finished!======================"   
        } || {
           echo ${mip}>>./log/failed.log
           echo "======================${mip} task failed!======================"
        }
        echo >&6
    } &
    sleep 1
done
wait 
exec 6>&-
exit 0

