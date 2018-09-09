#!/usr/bin/python 
# -*- coding:utf-8 -*-

import os
import re
import sys
import paramiko

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
    os.rename(outfile, infile) 
    


def ssh(ip, port, username, command):
    password=""
    ssh = paramiko.SSHClient()
    ssh.set_missing_host_key_policy(paramiko.AutoAddPolicy())
    ssh.connect("%s"%ip,port,"%s"%username, "%s"%password)
    stdin, stdout, stderr = ssh.exec_command("%s"%command)
    return stdout.readlines()
    ssh.close()

def doit(module, commands):
     if module == "all":
         r1 = re.compile(r' ')
     else:
         r1 = re.compile(r'%s'%module)
     f=open('ip_set','r') 
     arr = {}
     for i in f.readlines():
         if r1.search('%s' % i.strip()):
             ip = i.split()[0]
             remote=ssh("%s" % ip, 22, "root", "%s"%commands)
             arr[ip]=remote
     f.close()
     return arr

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
    ip_node=open('ip_set', 'r')
    allinfo=ip_node.readlines()
    for i in allinfo:
       node_groups.append(i.split()[1])
    node_groups.append('all')
    node_groups=list(set(node_groups))
    ip_node.close()

    if mod_type in node_groups:
        module_types = {
          '%s'%mod_type: doit("%s"%mod_type, commands),
          None: default,
        }
        result=module_types.get('%s'%mod_type)
        if result:
            print result
#            for key in result.keys():
#                print key
    else:
        print "No such group!"
        sys.exit()
