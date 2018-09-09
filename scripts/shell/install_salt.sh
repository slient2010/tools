#!/bin/bash
#
# Desc: salt客户端安装脚本(包含salt-master, salt-minion)
# Date: 2014-02-28
#

# 设置salt软件安装包下载地址和安装目录
download_url="www.qseeking.com:8080"
install_dir="/opt"

# 配置参数
master_ip="www.qseeking.com"


# 检测是否已经安装了salt
function check_install()
{
    if [ ! -d "${install_dir}/salt/" ]
    then
        echo "No salt install directory found, maybe you have not install it yet!"
        exit -1
    fi
}

# 安装或卸载salt
case $1 in
    install)
        # 安装salt 
        # 判断是否已经安装了salt
        if [ -d "${install_dir}/salt" ] 
        then
            echo "Already install salt" 
            exit -1
        fi
         
        # 创建临时目录,用于临时存放安装软件
        [ ! -d "/tmp/install" ] && mkdir /tmp/install/
        cd /tmp/install
        # 检查是否安装了xz压缩工具
        which xz > /dev/null 2>&1 
        # 若未安装则下载xz并解压安装
        if [ $? -ne 0 ] 
        then
            echo "Install xz tools first!"
            exit 1
            wget http://${download_url}/xz-5.0.5.tar.gz
            tar zxvf xz-5.0.5.tar.gz
            mv xz /usr/local/
            ln -s /usr/local/xz/bin/xz /usr/bin/
            ln -s /usr/local/xz/lib/liblzma.so.5.0.5 /usr/lib/liblzma.so.5
            ldconfig
        fi
        # 创建salt安装目录，由于编译时指定了目录，故必须使用指定目录来安装salt
        [ ! -d "${install_dir}" ] && mkdir ${install_dir}
        # 下载salt安装包    
        wget http://${download_url}/salt-0.17.5.tar.xz
        xz -d salt-0.17.5.tar.xz
        tar -xvf salt-0.17.5.tar
        mv salt ${install_dir}/
        # 创建链接，便可以不用输入绝对路径就可以使用salt相关命令
        need_lnk=(`ls ${install_dir}/salt/bin/ | grep salt | awk '{printf $1" "}'`)
        for lnk in ${need_lnk[@]}
        do
            ln -s ${install_dir}/salt/bin/${lnk} /usr/bin/
        done
        # salt使用到的动态库文件 
        ln -s ${install_dir}/salt/lib/libzmq.so.3.0.0 /usr/lib/libzmq.so.3
        ldconfig
        # 创建开机启动服务
        ln -s ${install_dir}/salt/salt-minion /etc/init.d/
        ln -s ${install_dir}/salt/salt-master /etc/init.d/
        chkconfig salt-minion on
        chkconfig salt-master on
        
        # 安装完成，删除临时文件
        rm /tmp/install/salt-0.17.5.tar
        [ -f "/tmp/install/xz-5.0.5.tar.gz" ] && rm /tmp/install/xz-5.0.5.tar.gz
        echo "Great, install salt finished!" 
    ;;

    configure)
        # salt的配置
        # 测试是否已经安装了salt，若没安装，则无需配置，直接退出程序
        check_install
        echo "`hostname`" > /opt/salt/etc/minion_id
        # 作为salt minion时候，需要配置/${install_dir}/etc/minion指向master的ip。
        sed -i 's/127.0.0.1/'${master_ip}'/' ${install_dir}/salt/etc/minion
    ;;
 
    uninstall)
        # 卸载安装好了的salt 
        # 测试是否已经安装了salt，若没安装，则无需卸载，直接退出程序
        check_install
        # 删除创建的salt快捷方式
        need_del_lnk=(`ls ${install_dir}/salt/bin/ | grep salt | awk '{printf $1" "}'`)
        for del_lnk in ${need_del_lnk[@]} 
        do
            rm /usr/bin/${del_lnk}
        done
        # 删除动态库文件链接
        rm  /usr/lib/libzmq.so.3

        # 删除启动脚本
        chkconfig salt-minion off
        chkconfig salt-master off
        rm /etc/init.d/salt-minion
        rm /etc/init.d/salt-master
        # 删除salt安装目录
        rm -rf ${install_dir}/salt/
        
        # 卸载完成
        echo "Uninstall salt finished!"
    ;;

    *)
        # 使用方法提示
        echo "Usage: ./$0 {install|configure|uninstall}"
    ;;
esac
