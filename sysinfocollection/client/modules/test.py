#!/usr/bin/env python
# -*- coding:utf-8 -*-
# data collector
# 数据收集客户端
# 时间：2014-10-21
#
'''
   扩展在此处实现
'''

from get_sys_info import *
import json

def sys_info():
#    sys_info = []
    disk_usage = get_disk_usage()
    load_info = get_load_info()
    mem_info = get_memory_info()
    sys_info="%s+%s+%s" %(disk_usage, load_info, mem_info)
#   sys_info.append(disk_usage)
#   sys_info.append(load_info)
#   sys_info.append(mem_info)
    return sys_info

if __name__ == "__main__":
    abc=sys_info()
    print type(abc)
#    print encodedjson


#def socket_start(remote_server, remote_port):
#    sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
#    sock.connect(("%s" % remote_server, int(remote_port)))
#    data = str(sys_info())
#    sock.sendall(data)
#    # 从服务端获取数据，暂时无用，注释
#   #result = sock.recv(1024)
#   #return result
#    sock.close()
