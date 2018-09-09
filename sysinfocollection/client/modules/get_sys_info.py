#!/usr/bin/env python
# -*- coding: utf-8 -*-
#
# Author : Jun
# Desc: 获取系统基本信息：1.挂载分区使用情况；2.系统负载；3.内存使用百分比；4.网卡实时流量 
# Date: 2014-10-16 
#


import os
import struct
import socket
import fcntl
from time import sleep

#用于获取系统已挂载分区的基本信息                                       
def get_disk_partitions(all=False):
    phydevs = []
    f = open("/proc/filesystems", "r")
    for line in f:
        if not line.startswith("nodev"):
            phydevs.append(line.strip())

    retlist = []
    f = open('/etc/mtab', "r")
    for line in f:
        if not all and line.startswith('none'):
            continue
        fields = line.split()
        device = fields[0]
        mountpoint = fields[1]
        fstype = fields[2]
        if not all and fstype not in phydevs:
            continue
        if device == 'none':
            device = ''
        retlist.append(mountpoint)
    return retlist

def get_disk_usage():
    disk_info = []
    for par in get_disk_partitions():
        st = os.statvfs(par)
#        free = int(st.f_bavail * st.f_frsize / 1000.0 / 1000.0 )
#        total= int(st.f_blocks * st.f_frsize / 1000.0 / 1000.0 )
        free = st.f_bfree * st.f_bsize / pow(1000.0,3)
        total = st.f_blocks * st.f_bsize / pow(1000.0,3)
        use_rate = (1 - free / total) * 100
        fmt_use_rate = "%.3f" % use_rate
#        stat = {'partition':par, 'disk_use_rate': fmt_use_rate}
        stat = {'%s'%par:fmt_use_rate}
        disk_info.append(stat)
#    disk_info=[{'disk_info':disk_info}]
    disk_info=[disk_info]
    return disk_info

#用于获取系统load信息 
def get_load_info():
    load_file = '/proc/loadavg'
    try:
        f = open(load_file)
    except:
        return []
    else:
        load_data =  f.readline()
        curr_load = load_data.split()
        return [{'load_status':curr_load[0]}]

#用于获取系统内存基本信息
def get_memory_info() :
    temp_mem = {}
    meminfo = '/proc/meminfo'
    try:
        f = open(meminfo)
    except:
        return []
    else:
        lines = f.readlines()
        f.close()

        for line in lines:
            if len(line) < 2: continue
            name = line.split(':')[0]
            var = line.split(':')[1].split()[0]
            temp_mem[name] = long(var) * 1024.0
       
        use_rate = (temp_mem['MemTotal']-temp_mem['MemFree']-temp_mem['Buffers']-temp_mem['Cached'])/temp_mem['MemTotal'] * 100
        fmt_use_rate = "%.3f" % use_rate
        return [{'mem_use_rate': fmt_use_rate}]


#用于获取系统中所有物理网卡设备的IP地址
def get_all_net_hardware():
    SIOCGIFCONF = 0x8912    
    ifreq_size = 24 + 2 * len( struct.pack( 'P', 0 ) )
    max_possible = 128
    bytes = max_possible * ifreq_size

    ifreq_buf = array.array( 'B', '\0' * bytes )
    ifconf = struct.pack( 'iP', bytes, ifreq_buf.buffer_info()[0] )

    s = socket.socket( socket.AF_INET, socket.SOCK_DGRAM )
    result = fcntl.ioctl( s.fileno(), SIOCGIFCONF, ifconf )
    s.close()

    ifc_len, ifreq = struct.unpack( 'iP', result )
    hd_list = []
    for i in range( 0, ifc_len, ifreq_size ):
        ifr_name = ( ifreq_buf.tostring()[ i : i + ifreq_size ] )[:16]
        dev_name = ifr_name.split( '\0', 1 )[0]
        hd_list.append(dev_name)
    return hd_list

def get_ip_from_hardware(dev_name="eth0"):
    SIOCGIFADDR = 0x8915
    ifreq = struct.pack( '16sH14s', dev_name, 0, '' )

    s = socket.socket( socket.AF_INET, socket.SOCK_DGRAM )
    ret = fcntl.ioctl( s.fileno(), SIOCGIFADDR, ifreq )
    s.close()

    addr = struct.unpack( '16sH14B', ret )[2:]
    return "%d.%d.%d.%d" %( addr[2], addr[3], addr[4], addr[5] )

def get_ip_info():
    hd_addr = {}
    for ifname in self.get_all_net_hardware():
       ipaddr = self.get_ip_from_hardware(ifname)
       hd_addr[ifname] = ipaddr
       
    return hd_addr

def get_hd_status():
    net_info = []
    start = _get_current_hd_status()
    sleep(1)
    end = _get_current_hd_status()
    
    for dev in start.keys():
        in_bw  = (end[dev][0] - start[dev][0]) / pow(1024.0,2)
        out_bw = (end[dev][2] - start[dev][2]) / pow(1024.0,2)
        
        fmt_in_bw = "%.3f" % (in_bw)
        fmt_out_bw = "%.3f" % (out_bw)
        intf = dict(
            zip(
                ( 'if_name','in_b','out_b'),
                ( dev, fmt_in_bw, fmt_out_bw)
                )
            )
    
        net_info.append(intf)
    
    return net_info
        
def _get_current_hd_status():
    net = {}
    f = open("/proc/net/dev")
    lines = f.readlines()
    f.close()
    for line in lines[2:]:
        dev_name, dev_status = line.split(":")
        if dev_name.strip() <> "lo":
            con = dev_status.split()
        net[dev_name.strip()] = (int(con[0]), int(con[1]), int(con[8]) ,int(con[9]))
    return net

