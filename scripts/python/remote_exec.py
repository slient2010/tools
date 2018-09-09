#!/usr/bin/python 
# -*- coding:utf-8 -*-

import os
import re
import sys
import thread
import paramiko
from multiprocessing.dummy import Pool as ThreadPool 

def delblankline(infile, outfile):
    """ Delete blanklines of infile """
    infp = open(infile, "r")
    outfp = open(outfile, "w")
    lines = infp.readlines()
    for li in lines:
        if li.split():
            outfp.writelines(li)
    infp.close()
    outfp.close()

def ssh(ips):
    print ips
    sys.exit()
    username="root"
    port=22
    password=""
    ip=ips[0].split(" ")[0] 
    command="%s"%ips[0].split(" ")[1]
    sshclient = paramiko.SSHClient()
    sshclient.set_missing_host_key_policy(paramiko.AutoAddPolicy())
    sshclient.connect("%s"%ip,port,"%s"%username, "%s"%password)
    stdin, stdout, stderr = sshclient.exec_command("%s"%command)
    result=stdout.readlines()
    results={}
    results={ip: result}
    sshclient.close()
    return results

def doit(module, commands):
    f=open('ip_sets','r') 
    info=f.readlines()
    r1 = re.compile(r'%s'%module)

    ips_info=[]
    for i in info:
        if r1.search('%s' % i.strip()):
            ip = i.split()[0]
            ips = '%s %s'%(ip,commands)
            ips_info.append(ips) 
    ips_info=transfer(ips_info)
    print ips_info
    pool = ThreadPool(4) 
    results = pool.map(ssh, ips_info)
    pool.close() 
    pool.join() 
    if results:
       return results
    f.close()

def transfer(ips_info):
    c = list(set(ips_info))
    b=[]
    for m in range(0, len(c)):
        d=[c[m]]
        b.append(d)
    return b

def default():
    return "default"

if __name__ == "__main__":
    if len(sys.argv) < 3:
        print "Usage: %s mod_type commands\nExample: %s dbank-10-node 'ls'" % (sys.argv[0], sys.argv[0])
        sys.exit()
    mod_type=sys.argv[1]
    commands=sys.argv[2]
    delblankline("ip_set", "ip_sets")
    node_groups=[]
    ip_node=open('ip_sets', 'r')
    allinfo=ip_node.readlines()
    for i in allinfo:
       node_groups.append(i.split()[1])
    node_groups.append('all')
    node_groups=list(set(node_groups))
    ip_node.close()

    if mod_type in node_groups:
        if mod_type != "all":
            module_types = {
              '%s'%mod_type: doit("%s"%mod_type, commands),
              None: default,
            }
            result=module_types.get('%s'%mod_type)
            if result:
                print result
        else:
            ips_info=[]
            for i in allinfo:
                ip = i.split()[0]
                ips = '%s %s'%(ip,commands)
                ips_info.append(ips) 
            ips_info=transfer(ips_info)
            print ips_info
            sys.exit()
#            pool = ThreadPool(4) 
#            results = pool.map(ssh, ips_info)
#            pool.close() 
#            pool.join() 
            results = ssh(ips_info)
            if results:
                print results
    else:
        print "No such group!"
        sys.exit()

