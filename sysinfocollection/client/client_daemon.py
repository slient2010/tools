#!/usr/bin/python
# -*-coding:utf-8-*-
#
#  FileName: agent_daemon.py
#  Author: Jun
#  Desc：日志收集守护进程
#  Date: 2014-10-17 17:00
#

import sys, os, time, atexit, string
from signal import SIGTERM 
import time
#import httplib
import commands
import ConfigParser
import codecs
from modules import mokylin_client


class Daemon:
    #def __init__(self, pidfile, stdin=stdin, stdout=stdout, stderr=stderr):
    def __init__(self, pidfile, stdin, stdout, stderr):
        self.pidfile= pidfile 
        self.stdin = stdin
        self.stdout = stdout
        self.stderr = stderr

    #创建守护子进程
    def _daemonize(self):
        try:
            pid = os.fork()
            # 退出主进程
            if pid > 0:
                sys.exit(0)
        except OSError, e:
            sys.stderr.wirte('fork #1 failed: %d (%s)\n' % (e.errno, e.strerror))
            sys.exit(1)
        os.chdir("/")
        os.umask(0)
        os.setsid()

        try:
            pid = os.fork()
            if pid > 0:
                sys.exit(0)
        except OSError, e:
            sys.stderr.wirte('fork #2 failed: %d (%s)\n' % (e.errno, e.strerror))
            sys.exit(1)

        #进程已经是守护进程，重定向标准文件描述符
        for f in sys.stdout, sys.stderr: f.flush()
        si = file(self.stdin, 'r')
        so = file(self.stdout,'a+')
        se = file(self.stderr,'a+',0)
        os.dup2(si.fileno(), sys.stdin.fileno())
        os.dup2(so.fileno(), sys.stdout.fileno())
        os.dup2(se.fileno(), sys.stderr.fileno())
            
        #创建processid文件
        atexit.register(self.delpid)
        pid = str(os.getpid())
        file(self.pidfile,'w+').write('%s\n' % pid)

    def delpid(self):
        os.remove(self.pidfile)

    def _run(self):
        # 服务开启时间
        self.start_time = time.strftime('%Y-%m-%d %H:%M:%S',time.localtime(time.time()))
        os.system("echo '[%s] Starting data_agent daemon...'" % self.start_time)
        # 循环执行，像服务端发送数据
        while True:
            try:
                # 启动信息收集脚本
                mokylin_client.socket_start(str(self.server), int(self.server_port))
                # 休眠sleeptime秒
                time.sleep(float(sleeptime))
            except OSError, err:
                print err


    #读取配置文件config.ini
    def readConfig(self, config_file):
        self.config_file = config_file
        cp =  ConfigParser.ConfigParser()
        try:
            cp.readfp(codecs.open(config_file, 'rw', 'utf-8'))
        #except IOError as e:
        except IOError, e:
            print "Open %s failed, no such file." % (config_file)
            #if raw_input('Press any key to continue'):
            sys.exit()

        self.server= cp.get('main', 'server')
        self.server_port = cp.get('main', 'port')
        self.role = cp.get('main', 'roles')
        self.services_num = cp.get('servers', 'running')
    
    def start(self, config):
        self.config = config
        #检查pid文件是否存在以探测是否存在进程
        try:
            pf = file(self.pidfile,'r')
            pid = int(pf.read().strip())
            pf.close()
        except IOError:
            pid = None
        
        if pid:
            message = 'pidfile %s already exist. Daemon already running?\n'
            sys.stderr.write(message % self.pidfile)
            sys.exit(1)
        
        #启动监控
        self.readConfig(self.config)
        self._daemonize()
        self._run()

    def stop(self):
        # 服务结束时间
        end_time = time.strftime('%Y-%m-%d %H:%M:%S',time.localtime(time.time()))
        #os.system("echo '[%s] Stop listening...'" % end_time)
        #读取pid
        try:
            pf = file(self.pidfile, 'r')
            pid = int(pf.read().strip())
            pf.close()
        except IOError:
            pid = None
     
        if not pid:
            message = 'pidfile %s does not exist. Daemon not running?\n'
            sys.stderr.write(message % self.pidfile)
            return #重启不报错

        #杀进程
        try:
            while 1:
                os.kill(pid, SIGTERM)
                time.sleep(0.1)
                file_obj = open(self.stdout, 'aw')
                file_obj.write("[%s] Stop listening all services....\n" % end_time)
                file_obj.close()
                print "[%s] Stop listenng all services...." % end_time
                # 调用系统命令，可执行相关命令
                # os.system('echo "test"')
        except OSError, err:
            err = str(err)
            if err.find('No such process') > 0:
                if os.path.exists(self.pidfile):
                    os.remove(self.pidfile)
            else:
                print str(err)
                sys.exit(1)
    
    def restart(self, config):
        self.config = config
        self.stop()
        self.start(self.config)

#程序入口
if __name__ == "__main__":

    if len(sys.argv) == 2:
        path_pre = os.getcwd()
        config= path_pre+ '/' + 'config.ini'

        # read configuration from file 
        if os.path.exists(config):
            #读取配置文件config.ini
            cp =  ConfigParser.ConfigParser()
            try:
                cp.readfp(codecs.open(config, 'r', 'utf-8'))
            #except IOError as e:
            except IOError, e:
                print "Open %s failed, no such file." % (config_file)
                #if raw_input('Press any key to continue'):
                sys.exit()
           
            server= cp.get('main', 'server')
            server_port = cp.get('main', 'port')
            sleeptime = cp.get('main', 'sleeptime')
           #role = cp.get('main', 'roles')
           #services_num = cp.get('servers', 'running')
            log_path = cp.get('logs', 'logpath')
            logname = cp.get('logs', 'logname')
            
            log_path=path_pre + '/' + log_path + '/' + logname
            daemon = Daemon('/tmp/watch_mokylin_client_process.pid', '/dev/null', '%s' % log_path, '%s' % log_path)
        else:
            print "配置文件不存在"
            sys.exit()

        if 'start' == sys.argv[1]:
            daemon.start('%s' % config)
        elif 'stop' == sys.argv[1]:
            daemon.stop()
        elif 'restart' == sys.argv[1]:
            daemon.restart('%s' % config)
        else:
            print 'Unknown command'
            sys.exit(2)
        sys.exit(0)
    else:
      print 'usage: %s (start|stop|restart)' % sys.argv[0]
      sys.exit(2)
