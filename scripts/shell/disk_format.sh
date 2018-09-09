#!/bin/bash

#鉴权，非root不能使用该脚本
if [ $UID -ne 0 ]
then
    echo "You have no rights to format disks!"
    exit
fi

#询问是否分区
echo 
echo "Do you really want to re-parted all disks?"
if read -t 10 -p "Yes/No: " commond
then
    case $commond in
        Yes)
            
            #disks=(`ls /sys/block/|grep sd |grep -v sda |awk '{printf " "$1}'`)
            disks=(`for i in sd{c..f}; do    echo $i; done`)
            for i in ${disks[@]}
            do
            {
                parted /dev/$i > /dev/null 2>&1 <<EOF
rm 1
rm 2
mklabel gpt
yes
mkpart primary 0 3001G
I
quit
EOF
            } &
            done
            wait 
            
            partprobe
            sleep 2
            
            for j in ${disks[@]}
            do
            {
                mkfs.reiserfs -f -q /dev/$j\1 >/dev/null
            } &
            done
            wait
            sleep 1
            partprobe
            echo "Finished!"
        ;;
        No|*)
            echo "Not parted and format any disks!"
            exit 
    esac
else
    echo "sorry, think it carefull before you make the decision!"
    exit
fi

