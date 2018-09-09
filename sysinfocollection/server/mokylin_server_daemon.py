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
import httplib
import commands
import ConfigParser
import codecs
import logging
from modules.check_server import *
from modules import mokylin_server 

class Daemon:
    #def __init__(self, pidfile, stdin=stdin, stdout=stdout, stderr=stderr):
    def __init__(self, pidfile, stdin, stdout, stderr, logger, dbconfig):
        self.pidfile= pidfile 
        self.stdin = stdin
        self.stdout = stdout
        self.stderr = stderr
        self.logger = logger
        self.dbconfig = dbconfig

    #创建守护子进程
    def _daemonize(self):
        try:
            pid = os.fork()
            # 退出主进程
            if pid > 0:
                sys.exit(0)
        except OSError, e:
            sys.stderr.wirte('fork #1 failed: %d (%s)\n' % (e.errno, e.strerror))
            logger.info('fork #1 failed: %d (%s)\n' % (e.errno, e.strerror))
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
            logger.info('fork #2 failed: %d (%s)\n' % (e.errno, e.strerror))
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
        file(self.pidfile,'aw+').write('%s\n' % pid)

    def delpid(self):
        os.remove(self.pidfile)

    def _run(self, dbconfig):
        self.dbconfig = dbconfig
        message = 'Starting data_server daemon...' 
        logger.info("%s" % message)
        # 开启监听服务
        mokylin_server.socket_start(str(self.server), int(self.server_port), self.dbconfig)
        while True:
            # 生成对象
            #check=check_server.checkService()
            check=checkService()
            # 检查进程结果
            check_result = check.check_service('mokylin_server.py')
            if int(check_result) == 0:
                # 启动信息收集脚本
                mokylin_server.socket_start(str(self.server), int(self.server_port), self.dbconfig)
                # 发送邮件
                check.sendmail('mokylin_server')
            # 休眠10秒
            # time.sleep(float(sleeptime))


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
        #self.services_num = cp.get('servers', 'running')
    
    def start(self, config, dbconfig):
        # 服务开启时间
        self.start_time = time.strftime('%Y-%m-%d %H:%M:%S',time.localtime(time.time()))
        # os.system("echo '[%s] Starting data_agent daemon...'" % self.start_time)
        print "[%s] Starting listening data_server services...." % self.start_time
        # 日志记录
        self.config = config
        # 数据库配置
        self.dbconfig = dbconfig
        #检查pid文件是否存在以探测是否存在进程
        try:
            pf = file(self.pidfile,'r')
            pid = int(pf.read().strip())
            pf.close()
        except IOError:
            pid = None
        
        if pid:
            message = 'pidfile %s already exist. Daemon already running?\n'
            logger.info("%s" % message)
            sys.stderr.write(message % self.pidfile)
            sys.exit(1)
        
        #启动监控
        self.readConfig(self.config)
        self._daemonize()
        self._run(self.dbconfig)

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
            logger.info("%s" % message)
            sys.stderr.write(message % self.pidfile)
            return #重启不报错

        #杀进程
        try:
            while 1:
                os.kill(pid, SIGTERM)
                time.sleep(0.1)
                #file_obj = open(self.stdout, 'aw')
                # file_obj.write("[%s] INFO Stop listening all services....\n" % end_time)
                #file_obj.close()
                message = "Stop listening all services...."
                print "[%s] Stop listening all services...." % end_time
                logger.info("%s" % message)
                # 调用系统命令，可执行相关命令
                # os.system('echo "test"')
        except OSError, err:
            err = str(err)
            if err.find('No such process') > 0:
                if os.path.exists(self.pidfile):
                    os.remove(self.pidfile)
            else:
                print str(err)
                logger.info("%s" % str(err))
                sys.exit(1)
    
    def restart(self, config, dbconfig):
        self.config = config
        self.dbconfig = dbconfig
        self.stop()
        self.start(self.config, self.dbconfig)



#程序入口

if __name__ == "__main__":
    # 配置文件处理，这里有个问题，异常告警处理，暂时忽略
    if len(sys.argv) == 2:
        path_pre = os.getcwd()
        config= path_pre+ '/' + 'config.ini'
        #读取配置文件config.ini
        if os.path.exists(config):
            cp =  ConfigParser.ConfigParser()
            try:
                cp.readfp(codecs.open(config, 'r', 'utf-8'))
            # 读取日志异常
            except IOError, e:
                print "Open %s failed, no such file." % (config)
                sys.exit()
            # 监听地址 
            server= cp.get('main', 'server')
            # 监听端口
            server_port = cp.get('main', 'port')
            # 日志路径
            log_path = cp.get('logs', 'logpath')
            # 日志名称
            logname = cp.get('logs', 'logname')
            # 日志级别
            loglevel = cp.get('logs', 'loglevel')
            # 日志绝对路径 
            log_path = path_pre + '/' + log_path + '/' + logname
            dbconfig = []
            # 数据库地址
            dbaddress = cp.get('Database', 'ip')
            dbaddress = {'dbaddress': dbaddress}
            dbconfig.append(dbaddress)
            # 数据库端口
            dbport= cp.get('Database', 'port')
            dbport = {'dbport': dbport}
            dbconfig.append(dbport)
            # 数据库用户名
            dbuser= cp.get('Database', 'user')
            dbuser = {'dbuser': dbuser}
            dbconfig.append(dbuser)
            # 数据库密码
            dbpass = cp.get('Database', 'password')
            dbpass = {'dbpass': dbpass}
            dbconfig.append(dbpass)
            # 数据库库名
            dbname = cp.get('Database', 'dbname')
            dbname = {'dbname': dbname}
            dbconfig.append(dbname)


            # 日志格式
            logging.basicConfig(level=str(loglevel),
                    format='[%(asctime)s] %(levelname)-4s %(message)s',
                    datefmt='%Y-%m-%d %H:%M:%S',
                    filename='%s' % log_path,
                    filemode='aw')
            logger = logging.getLogger("logstart")
            daemon = Daemon('/tmp/watch_process.pid', '/dev/null', '%s' % log_path, '%s' % log_path, logger, dbconfig)
        else:
            print "配置文件不存在"
            sys.exit()

        if 'start' == sys.argv[1]:
            daemon.start('%s' % config, dbconfig)
        elif 'stop' == sys.argv[1]:
            daemon.stop()
        elif 'restart' == sys.argv[1]:
            daemon.restart('%s' % config, dbconfig)
        else:
            print 'Unknown command'
            sys.exit(2)
        sys.exit(0)
    else:
      print 'usage: %s (start|stop|restart)' % sys.argv[0]
      sys.exit(2)

